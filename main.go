package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"semtag/pkg/docker"
	"semtag/pkg/git"
	"semtag/pkg/version"
)

var (
	dryRun  bool
	skipInc bool
	suffix  string
	prefix  string
)

var (
	dockerImage      string
	dockerRepository string
	shouldTagGit     bool
	file             string
	fileVerPattern   string
)

func parseFlags() {
	flag.BoolVar(&dryRun, "dry-run", false, "if set, only print the object(s) that would be sent, without sending the data")
	flag.BoolVar(&skipInc, "skip-increment", false, "if set, do not increment the version number")
	flag.StringVar(&suffix, "suffix", "", `if set, append the suffix to the version number (e.g. "0.1.0-rc")`)
	flag.StringVar(&prefix, "prefix", "", `if set, append the prefix to the version number (e.g. "api-0.1.0")`)

	flag.StringVar(&dockerImage, "docker-image", "", "a Docker image saved as a tar archive (e.g. use `api.tar` for an image saved with `docker save api:latest > api.tar`)")
	flag.StringVar(&dockerRepository, "docker-repository", "", "target Docker repository (e.g. '$MY_DOCKER_REGISTRY/$MY_APP_NAME')")

	flag.BoolVar(&shouldTagGit, "git-tag", false, "if set, create an annotated tag")

	flag.StringVar(&file, "file", "", `a file that contains the version number (e.g. "setup.py")`)
	flag.StringVar(&fileVerPattern, "file-version-pattern", "", `the pattern expected for the file version (e.g. "version='%s',")`)

	flag.Parse()
	if len(os.Args) == 1 {
		log.Println("Please see usage:")
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
	ver, nextVer := getVersions()

	if len(dockerImage) > 0 && len(dockerRepository) > 0 {
		tagDocker(ver)
	}

	if shouldTagGit {
		tagGit(nextVer)
	}

	if len(file) > 0 && len(fileVerPattern) > 0 {
		tagFile(nextVer)
	}
}

func getVersions() (*version.Version, *version.Version) {
	var v, nextV version.Version
	v.Suffix = suffix
	v.Prefix = prefix
	v = *v.GetLatest()
	fmt.Println("version:", v.String())
	if skipInc {
		log.Println("skip version increment: flag set by user")
		return &v, &v
	}
	nextV = *v.GetLatest()
	nextV.IncrementAuto()
	fmt.Println("next-version:", nextV.String())
	return &v, &nextV
}

func tagDocker(ver *version.Version) {
	docker.Load(dockerImage)
	img := &docker.Image{
		Name:                strings.Replace(dockerImage, ".tar", "", 1),
		Tags:                ver.AsList(git.DescribeLong()),
		ContainerRepository: dockerRepository,
	}
	log.Println("tag docker image:", img)
	if !dryRun {
		img.Tag()
		img.Push()
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

func tagFile(ver *version.Version) {
	const commitMsgVerBump = "(v_bmp)"
	f := version.File{
		Path:          file,
		VersionFormat: fileVerPattern,
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
		git.Add(file)
		msg, err := git.GetLastCommitNames(-1)
		if err != nil {
			log.Fatal(err)
		}
		git.Commit(commitMsgVerBump + ": " + *msg)
		git.Push("--tags")
	}
}
