package vethhunter

import (
	"fmt"
	"strings"
	"syscall"

	"github.com/Sirupsen/logrus"
	"github.com/fsouza/go-dockerclient"
	"github.com/vishvananda/netlink"
	"github.com/vishvananda/netns"
)

// VethHunter holds the info of docker client, etc.
type VethHunter struct {
	DC *docker.Client
}

// NewVethHunterFromLocalDocker returns a new hunter
// for local docker socket
func NewVethHunterFromLocalDocker() *VethHunter {

	endpoint := "unix:///var/run/docker.sock"
	client, _ := docker.NewClient(endpoint)

	return &VethHunter{DC: client}
}

// GetHostVethOfContainer gets the host side of the veth
func (vh *VethHunter) GetHostVethOfContainer(cid string) (string, error) {
	c, err := vh.DC.InspectContainer(cid)
	if err != nil {
		return "", err
	}

	nsPath := c.NetworkSettings.SandboxKey
	if strings.Contains(nsPath, "default") {
		return "", nil
	}

	nsh, err := netns.GetFromPath(nsPath)
	if err != nil {
		return "", err
	}

	nlh, err := netlink.NewHandleAt(nsh, syscall.NETLINK_ROUTE)
	if err != nil {
		return "", err
	}

	hostnlh, err := netlink.NewHandleAt(netns.None(), syscall.NETLINK_ROUTE)
	if err != nil {
		return "", err
	}

	links, err := nlh.LinkList()
	if err != nil {
		logrus.Errorf("failed to list interfaces: %v", err)
	}

	for _, l := range links {
		if l.Type() == "veth" {
			hostPeer, err := hostnlh.LinkByIndex(l.Attrs().Index + 1)
			if err != nil {
				return "", err
			}
			hostVeth := hostPeer.Attrs().Name
			return hostVeth, nil
		}
	}

	return "", fmt.Errorf("Invalid container id")
}
