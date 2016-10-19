package main

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/fsouza/go-dockerclient"
	"github.com/leodotcloud/vethhunter/vethhunter"
)

func init() {
	//logrus.SetLevel(logrus.DebugLevel)
}

func main() {
	endpoint := "unix:///var/run/docker.sock"
	client, _ := docker.NewClient(endpoint)

	containers, _ := client.ListContainers(docker.ListContainersOptions{})

	vh := &vethhunter.VethHunter{client}

	for _, c := range containers {
		//logrus.Debugf("c: %+v", c)
		hostVeth, err := vh.GetHostVethOfContainer(c.ID)
		if err != nil {
			logrus.Errorf("Error: %v", err)
		}
		fmt.Println(c.ID, ":", hostVeth)
	}
}
