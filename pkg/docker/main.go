package docker

import (
	"log"

	"semtag/pkg"
)

func Load(tarFile string) {
	out, err := pkg.Shell("docker load < " + tarFile)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(out)
}

func Tag(image string, remoteImage string) {
	out, err := pkg.Shellf("docker tag %s %s", image, remoteImage)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(out)
}

func Push(image string) {
	out, err := pkg.Shell("docker push " + image)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(out)
}
