package commands

import (
	"fmt"
	"os"

	"github.com/cheggaaa/pb"
	"github.com/simongui/cloudraker/containers"
	"github.com/simongui/cloudraker/mysql"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

// AddCommand Context for "add" command
type AddCommand struct {
	Cluster    *string
	Datacenter *string
	Host       *string
	IPAddress  *string
	host       string
	serverID   string
}

type addStepFunc func(cmd *AddCommand) error

// Run Executes the command like action.
func (cmd *AddCommand) Run(c *kingpin.ParseContext) error {
	steps := []addStepFunc{
		containerExistsStep,
		addContainerStep,
		addToDockerNetworkStep,
		addIPAliasStep,
		addDNSStep,
		startSSHStep,
		getMySQLResults,
	}

	progress := pb.New(len(steps)).Prefix(*cmd.Host)
	progress.ShowTimeLeft = true
	progress.ShowFinalTime = true

	pool, err := pb.StartPool(progress)
	if err != nil {
		panic(err)
	}

	for i, step := range steps {
		err := step(cmd)
		if err != nil {
			fmt.Printf("Failed step %d\n%s", i, err)
			pool.Stop()
			fmt.Println(err)
			os.Exit(1)
		}
		progress.Increment()
	}
	pool.Stop()

	fmt.Printf("MySQL running \n\thost: %s\n\tid: %s", cmd.host, cmd.serverID)

	return nil
}

func containerExistsStep(cmd *AddCommand) error {
	// Check if the node is already added.
	exists := containers.Exists(*cmd.Host)
	if exists {
		return fmt.Errorf("%s is already a member of an existing cluster", *cmd.Host)
	}
	return nil
}

func addContainerStep(cmd *AddCommand) error {
	err := containers.Add(*cmd.Cluster, *cmd.Datacenter, *cmd.Host, *cmd.IPAddress)
	return err
}

func addToDockerNetworkStep(cmd *AddCommand) error {
	containers.AddToDockerNetwork(*cmd.Cluster, *cmd.Host)
	return nil
}

func addIPAliasStep(cmd *AddCommand) error {
	containers.AddIPAlias(*cmd.IPAddress)
	return nil
}

func addDNSStep(cmd *AddCommand) error {
	containers.AddDNS(*cmd.Host, *cmd.IPAddress)
	return nil
}

func startSSHStep(cmd *AddCommand) error {
	containers.StartSSH(*cmd.Host)
	return nil
}

func getMySQLResults(cmd *AddCommand) error {
	var err error
	cmd.host, err = mysql.GetHostname("shard0-db1.local1.com", 33301, "root", "password", "30s")
	if err != nil {
		return err
	}
	cmd.serverID, err = mysql.GetServerID("shard0-db1.local1.com", 33301, "root", "password", "30s")
	if err != nil {
		return err
	}
	return nil
}
