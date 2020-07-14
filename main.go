package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"semtag/pkg"
	"semtag/pkg/docker"
	"semtag/pkg/git"
	"semtag/pkg/version"
)

var (
	push             bool
	incrementVersion bool
	suffix           string
	prefix           string
)

var (
	dockerImage      string
	dockerRepository string
	shouldTagGit     bool
	file             string
	fileVerPattern   string
)

func parseFlags() {
	flag.BoolVar(&push, "push", false, "if set, push the object(s)")
	flag.BoolVar(&incrementVersion, "increment", false, "if set, increment the version number")
	flag.StringVar(&suffix, "suffix", "", `if set, append the suffix to the version number (e.g. "0.1.0-rc")`)
	flag.StringVar(&prefix, "prefix", "", `if set, append the prefix to the version number (e.g. "api-0.1.0")`)

	flag.StringVar(&dockerImage, "docker-image", "", "a Docker image saved as a tar archive (e.g. use `api.tar` for an image saved with `docker save api:latest > api.tar`)")
	flag.StringVar(&dockerRepository, "docker-repository", "", "target Docker repository (e.g. '$MY_DOCKER_REGISTRY/$MY_APP_NAME')")

	flag.BoolVar(&shouldTagGit, "git-tag", false, "if set, create an annotated tag")

	flag.StringVar(&file, "file", "", `a file that contains the version number (e.g. "setup.py")`)
	flag.StringVar(&fileVerPattern, "file-version-pattern", "%s", `the pattern expected for the file version (e.g. "version='%s',")`)

	flag.Parse()
	fmt.Println(getVersion().String())
	if push {
		log.Println("push mode enabled")
	}
}

func main() {
	parseFlags()

	v := getVersion()

	if len(dockerImage) > 0 && len(dockerRepository) > 0 {
		tagDocker(v)
	}

	if shouldTagGit {
		tagGit(v)
	}

	if len(file) > 0 && len(fileVerPattern) > 0 {
		tagFile(v)
	}
}

func getVersion() *version.Version {
	var v version.Version
	v.Suffix = suffix
	v.Prefix = prefix
	v = *v.GetLatest()
	if incrementVersion {
		v.IncrementAuto()
	}
	return &v
}

func tagDocker(ver *version.Version) {
	docker.Load(dockerImage)
	img := &docker.Image{
		Name:                strings.Replace(dockerImage, ".tar", "", 1),
		Tags:                ver.AsList(git.DescribeLong()),
		ContainerRepository: dockerRepository,
	}
	log.Println("tag docker image:", img)
	img.Tag()
	if push {
		img.Push()
	}
}

func tagGit(ver *version.Version) {
	tag := &git.TagObj{
		Name: ver.String(),
	}
	tag.SetMessage()
	log.Println("tag git:", tag)
	git.SetGitConfig()
	if push {
		tag.Push()
	}
}

func tagFile(ver *version.Version) {
	const commitMsgVerBump = "chore(version): "
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
		log.Println("skip version increment: already incremented")
		return
	}
	newContents := f.ReplaceSubstring()

	log.Println("tag file:", f)
	if pkg.DEBUG != "" {
		log.Println("\n", *newContents)
	}
	f.Write(newContents)
	git.Add(file)
	if push {
		msg, err := git.GetLastCommitNames(-1)
		if err != nil {
			log.Fatal(err)
		}
		git.Commit(commitMsgVerBump + ": " + *msg)
		git.Push("--tags")
	}
}
