package main

import (
	"flag"
	"log"
	"os"
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
	flag.BoolVar(&skipInc, "skip-increment", false, "if true, do not increment the version number")

	flag.Parse()
	if tag == "" {
		log.Println("Parameter `tag` not found. Please see usage:")
		flag.PrintDefaults()
		os.Exit(101)
	}
	if dryRun {
		log.Println("dry run mode enabled")
	}
}

func main() {
	parseFlags()

	git.Fetch()
	_, nextVer := getVersions()

	switch tag {
	case "file":
		tagFile(nextVer)
	case "git":
		tagGit(nextVer)
	case "docker":
		tagDocker(nextVer)
	}
}

func getVersions() (*version.Version, *version.Version) {
	var v, nextV version.Version
	v.Suffix = suffix
	v.Prefix = prefix
	v = *v.GetLatest()
	log.Println("version:", v.String())
	if skipInc {
		log.Println("skip version increment: flag set by user")
		return &v, &v
	}
	nextV = *v.GetLatest()
	nextV.IncrementAuto()
	log.Println("next version:", nextV.String())
	return &v, &nextV
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
		git.Push("--tags")
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
		Tags:                ver.AsList(git.DescribeLong()),
		ContainerRepository: out,
	}
	log.Println("tag docker image:", img)
	if !dryRun {
		img.Tag()
		img.Push()
	}
}
