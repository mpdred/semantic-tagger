package docker

import (
	"fmt"
)

type Image struct {
	Name                string
	Tags                []string
	ContainerRepository string
}

func (img *Image) getRemoteNames() *[]string {
	var remoteNames []string
	for _, ver := range img.Tags {
		rName := fmt.Sprintf("%s:%s", img.ContainerRepository, ver)
		remoteNames = append(remoteNames, rName)
	}
	rNameLatest := fmt.Sprintf("%s:%s-%s", img.ContainerRepository, "latest", img.Name)
	remoteNames = append(remoteNames, rNameLatest)
	return &remoteNames
}

func (img *Image) Tag() {
	for _, rName := range *img.getRemoteNames() {
		Tag(img.Name, rName)
	}
}

func (img *Image) Push() {
	for _, rName := range *img.getRemoteNames() {
		Push(rName)
	}
}
