package main

import (
	"flag"
	"fmt"
	"log"

	"semtag/pkg/docker"
	"semtag/pkg/git"
	"semtag/pkg/version"
)

var (
	tag    string
	in     string
	out    string
	suffix string
	prefix string
	dryRun bool
)

func parseFlags() {
	flag.StringVar(&tag, "tag", "", "tag choices: file | git | docker")
	flag.StringVar(&in, "in", "", `input data: can be either 1) a docker image .tar file without the file extension (e.g. "api") or 2) a file that contains the version number (e.g. "setup.py")`)
	flag.StringVar(&out, "out", "", `output: can be either 1) a docker repository or 2) the pattern for the file version (e.g. "version='%s',")`)
	flag.StringVar(&suffix, "suffix", "", `if set, append the suffix to the version number (e.g. "0.1.0-rc")`)
	flag.StringVar(&prefix, "prefix", "", `if set, append the prefix to the version number (e.g. "api-0.1.0")`)
	flag.BoolVar(&dryRun, "dry-run", false, "if true, only print the object(s) that would be sent, without sending the data")
	flag.Parse()

	if dryRun {
		log.Println("dry run mode enabled")
	}
}

func main() {
	parseFlags()
	git.Fetch()

	var ver, nextVer version.Version
	ver.Suffix = suffix
	ver.Prefix = prefix
	ver = *ver.GetLatest()
	log.Println("current version:", ver.String())
	nextVer = *ver.GetLatest()
	changeType := nextVer.IncrementAuto()
	log.Println("next version:", nextVer.String())

	switch tag {
	case "file":
		f := version.File{
			Path:          in,
			VersionFormat: out,
			Version:       nextVer.String(),
		}
		commitMsg, err := git.GetLastCommitNames(-1)
		if err != nil {
			log.Fatal(err)
		}
		newContents := f.ReplaceSubstring()
		log.Println("new file contents\n", *newContents)
		if !dryRun {
			f.Write(newContents)
			git.Add(in)
			git.Commit(fmt.Sprintf("%s ver inc: %s %s", *commitMsg, nextVer.String(), changeType.String()))
			git.Push("")
		}

	case "git":
		tag := &git.TagObj{
			Name: nextVer.String(),
		}
		tag.SetMessage()
		log.Println("new git tag:", tag)
		if !dryRun {
			tag.Push()
		}

	case "docker":
		docker.Load(in + ".tar")
		img := &docker.Image{
			Name:                in,
			Tags:                ver.AsList(),
			ContainerRepository: out,
		}
		log.Println("new docker image tags", img)
		if !dryRun {
			img.Tag()
			img.Push()
		}
	default:
		log.Printf("not implemented for %q\n", tag)
		flag.PrintDefaults()
	}
}
