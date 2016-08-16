package containers

import (
	"errors"
	"fmt"
	"os/exec"
	"runtime"

	"golang.org/x/net/context"

	"github.com/docker/engine-api/client"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/filters"
)

// DOCKERLABELNAMESPACE Represents namespace all labels get stored.
const DOCKERLABELNAMESPACE = "cloudraker."

// DOCKERLABELCLUSTER Represents the cluster label name stored in Docker.
const DOCKERLABELCLUSTER = DOCKERLABELNAMESPACE + "cluster"

// DOCKERLABELDATACENTER Represents the datacenter label name stored in Docker.
const DOCKERLABELDATACENTER = DOCKERLABELNAMESPACE + "datacenter"

// DOCKERLABELIPADDRESS Represents the ip address label stored in Docker.
const DOCKERLABELIPADDRESS = DOCKERLABELNAMESPACE + "ipaddress"

// executeCommand Executes a shell command and returns the resulting output.
func executeCommand(command string, arguments []string) ([]byte, error) {
	var cmdOut []byte
	var err error

	if cmdOut, err = exec.Command(command, arguments...).Output(); err != nil {
		return cmdOut, err
	}
	return cmdOut, nil
}

// AddFilter Adds a filter to the Docker query.
func AddFilter(filter filters.Args, name string, value string) {
	filter.Add(name, value)
}

// AddLabelFilter Adds a filter to the Docker query.
func AddLabelFilter(filter filters.Args, name string, value string) {
	label := fmt.Sprintf("%s%s=%s", DOCKERLABELNAMESPACE, name, value)
	filter.Add("label", label)
}

// GetDockerContainers Returns all the docker containers.
func GetDockerContainers(filter filters.Args) ([]types.Container, error) {
	defaultHeaders := map[string]string{"User-Agent": "engine-api-cli-1.0"}
	cli, err := client.NewClient("unix:///var/run/docker.sock", "v1.22", nil, defaultHeaders)
	if err != nil {
		return nil, err
	}

	options := types.ContainerListOptions{All: true}
	options.Filter = filter

	containers, err := cli.ContainerList(context.Background(), options)
	if err != nil {
		panic(err)
	}
	return containers, nil
}

// GetDockerContainer Returns the container associated with the provided host.
func GetDockerContainer(host string) (*types.Container, error) {
	var err error
	filter := filters.NewArgs()
	containers := []types.Container{}

	AddFilter(filter, "name", host)
	containers, err = GetDockerContainers(filter)

	if err != nil {
		return &types.Container{}, err
	}
	if len(containers) > 0 {
		return &containers[0], nil
	}
	return &types.Container{}, nil
}

// Exists Returns whether the specified container exists.
func Exists(host string) bool {
	var err error
	filter := filters.NewArgs()
	containers := []types.Container{}

	AddFilter(filter, "name", host)
	containers, err = GetDockerContainers(filter)
	if err != nil {
		panic(err)
	}
	if len(containers) > 0 {
		return true
	}
	return false
}

// GetProcesses Returns list of processes running in the specified container.
func GetProcesses(host string) (types.ContainerProcessList, error) {
	var processes types.ContainerProcessList

	defaultHeaders := map[string]string{"User-Agent": "engine-api-cli-1.0"}
	cli, err := client.NewClient("unix:///var/run/docker.sock", "v1.22", nil, defaultHeaders)
	if err != nil {
		return processes, err
	}

	processes, err = cli.ContainerTop(context.Background(), host, nil)
	if err != nil {
		panic(err)
	}
	return processes, nil
}

// GetHasProcess Returns whether the given host has a process running.
func GetHasProcess(host string, processName string) (bool, error) {
	processes, err := GetProcesses(host)
	var hasProcess = false
	for _, process := range processes.Processes {
		if process[3] == processName {
			hasProcess = true
			break
		}
	}
	return hasProcess, err
}

// Add Adds a new container to the specified cluster and datacenter.
func Add(cluster string, datacenter string, host string, ipaddress string, mysqlPort int, sshPort int) error {
	args := []string{
		"run",
		"--name",
		host,
		"--hostname",
		host,
		"-e",
		"MYSQL_ROOT_PASSWORD=password",
		"-p",
		fmt.Sprintf("%d:%d", mysqlPort, 3306),
		"-p",
		fmt.Sprintf("%d:%d", sshPort, 22),
		"--label",
		fmt.Sprintf("%s=%s", DOCKERLABELCLUSTER, cluster),
		"--label",
		fmt.Sprintf("%s=%s", DOCKERLABELDATACENTER, datacenter),
		"--label",
		fmt.Sprintf("%s=%s", DOCKERLABELIPADDRESS, ipaddress),
		"-d",
		"cloudraker/mysql-server:5.7",
		"--server-id=1",
	}
	_, err := executeCommand("docker", args)

	if err != nil {
		fmt.Println(args)
		fmt.Println(err)
		return err
	}
	return nil
}

// StartSSH Starts the SSH daemon on the specified host.
func StartSSH(host string) {
	var running = false
	for ok := true; ok; ok = (running == false) {
		executeCommand("docker", []string{
			"exec",
			host,
			"service",
			"ssh",
			"start",
		})
		running, _ = GetHasProcess(host, "/usr/sbin/sshd")
	}
}

// AddToDockerNetwork Adds the specified host to the specified docker network.
func AddToDockerNetwork(network string, host string) {
	executeCommand("docker", []string{
		"network",
		"create",
		network,
	})
	executeCommand("docker", []string{
		"network",
		"connect",
		network,
		host,
	})
}

// RemoveDockerNetwork Removes the specified docker network.
func RemoveDockerNetwork(network string) {
	executeCommand("docker", []string{
		"network",
		"rm",
		network,
	})
}

// AddIPAlias Adds an alias for the specified IP address to the loopback adapter.
func AddIPAlias(ipaddress string) {
	if runtime.GOOS == "linux" {
		executeCommand("sudo", []string{
			"ifconfig",
			"lo:0",
			ipaddress,
			"up",
		})
	} else if runtime.GOOS == "darwin" {
		executeCommand("sudo", []string{
			"ifconfig",
			"lo0",
			ipaddress,
			"alias",
		})
	}
}

// RemoveIPAlias Removes the specified IP address alias from the loopback adapter.
func RemoveIPAlias(ipaddress string) {
	if runtime.GOOS == "linux" {
		executeCommand("sudo", []string{
			"ifconfig",
			"lo:0",
			ipaddress,
			"down",
		})
	} else if runtime.GOOS == "darwin" {
		executeCommand("sudo", []string{
			"ifconfig",
			"lo0",
			ipaddress,
			"delete",
		})
	}
}

// AddDNS Adds the specified dns entry pointing to the specified IP address.
func AddDNS(host string, ipaddress string) {
	executeCommand("sudo", []string{
		"ghost",
		"add",
		host,
		ipaddress,
	})
}

// RemoveDNS Removes the specified dns entry pointing to the specified IP address.
func RemoveDNS(host string) {
	executeCommand("sudo", []string{
		"ghost",
		"delete",
		host,
	})
}

// Remove Removes a new container from the specified cluster.
func Remove(container *types.Container) error {
	defaultHeaders := map[string]string{"User-Agent": "engine-api-cli-1.0"}
	cli, err := client.NewClient("unix:///var/run/docker.sock", "v1.22", nil, defaultHeaders)
	if err != nil {
		return err
	}

	options := types.ContainerRemoveOptions{ContainerID: container.ID, Force: true, RemoveLinks: false, RemoveVolumes: true}
	err = cli.ContainerRemove(context.Background(), options)
	if err != nil {
		panic(err)
	}

	// network := container.Labels[DOCKERLABELCLUSTER]
	// ipaddress := container.Labels[DOCKERLABELIPADDRESS]
	// RemoveDockerNetwork(network)
	// RemoveIPAlias(ipaddress)
	return nil
}

// RemoveHostsIPAlias Removes the specified hosts IP alias.
func RemoveHostsIPAlias(container *types.Container) error {
	if container.ID == "" {
		return errors.New("Unable to find container")
	}

	ipaddress := container.Labels[DOCKERLABELIPADDRESS]
	RemoveIPAlias(ipaddress)
	return nil
}
