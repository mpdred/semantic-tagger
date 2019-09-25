package main

import (
	"flag"
	"log"
	"strings"

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
)
var (
	dryRun  bool
	skipInc bool
)

const commitMsgVerBump = "(v_bmp)"

func parseFlags() {
	flag.StringVar(&tag, "tag", "", "tag choices: file | git | docker")
	flag.StringVar(&in, "in", "", `input data: can be either 1) a docker image .tar file without the file extension (e.g. "api") or 2) a file that contains the version number (e.g. "setup.py")`)
	flag.StringVar(&out, "out", "", `output: can be either 1) a docker repository or 2) the pattern for the file version (e.g. "version='%s',")`)
	flag.StringVar(&suffix, "suffix", "", `if set, append the suffix to the version number (e.g. "0.1.0-rc")`)
	flag.StringVar(&prefix, "prefix", "", `if set, append the prefix to the version number (e.g. "api-0.1.0")`)

	flag.BoolVar(&dryRun, "dry-run", false, "if true, only print the object(s) that would be sent, without sending the data")
	flag.BoolVar(&skipInc, "skip-inc", false, "if true, do not increment the version number")
	flag.Parse()

	if dryRun {
		log.Println("dry run mode enabled")
	}
}

func main() {
	parseFlags()
	git.Fetch()

	ver, nextVer := getVersions()

	switch tag {
	case "file":
		tagFile(nextVer)
	case "git":
		tagGit(nextVer)
	case "docker":
		tagDocker(ver)
	default:
		flag.PrintDefaults()
	}
}

func getVersions() (*version.Version, *version.Version) {
	var ver, nextVer version.Version
	ver.Suffix = suffix
	ver.Prefix = prefix
	ver = *ver.GetLatest()
	log.Println("current version:", ver.String())
	if skipInc {
		log.Println("skip version increment: flag set by user")
		return &ver, &ver
	}
	nextVer = *ver.GetLatest()
	nextVer.IncrementAuto()
	log.Println("next version:", nextVer.String())
	return &ver, &nextVer
}

func tagFile(ver *version.Version) {
	f := version.File{
		Path:          in,
		VersionFormat: out,
		Version:       ver.String(),
	}
	out, err := git.GetLastCommitNames(-1)
	if err != nil {
		log.Fatal(err)
	}
	if strings.Contains(*out, commitMsgVerBump) {
		log.Fatal("skip version increment: already incremented")
	}
	newContents := f.ReplaceSubstring()

	log.Println("tag file:", f, "\n", *newContents)
	f.Write(newContents)
	if !dryRun {
		git.Add(in)
		msg, err := git.GetLastCommitNames(-1)
		if err != nil {
			log.Fatal(err)
		}
		git.Commit(commitMsgVerBump + ": " + *msg)
		git.Push()
	}
}

func tagGit(ver *version.Version) {
	tag := &git.TagObj{
		Name: ver.String(),
	}
	tag.SetMessage()
	log.Println("tag git:", tag)
	if !dryRun {
		tag.Push()
	}
}

func tagDocker(ver *version.Version) {
	docker.Load(in + ".tar")
	img := &docker.Image{
		Name:                in,
		Tags:                ver.AsList(),
		ContainerRepository: out,
	}
	log.Println("tag docker image:", img)
	if !dryRun {
		img.Tag()
		img.Push()
	}
}
