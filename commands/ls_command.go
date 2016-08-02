package commands

import (
	"fmt"
	"os"
	"strconv"

	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/filters"
	"github.com/olekukonko/tablewriter"
	"github.com/simongui/cloudraker/containers"
	"gopkg.in/alecthomas/kingpin.v2"
)

// LsCommand Context for "ls" command
type LsCommand struct {
	Cluster    *string
	Datacenter *string
	Host       *string
}

// Run Executes the command like action.
func (cmd *LsCommand) Run(c *kingpin.ParseContext) error {
	var err error
	filter := filters.NewArgs()
	tableData := [][]string{}
	runningContainers := []types.Container{}

	if *cmd.Host != "" {
		// A specific host is provided so only list the details of that node.
		containers.AddFilter(filter, "name", *cmd.Host)
	} else if *cmd.Cluster != "" {
		containers.AddLabelFilter(filter, "cluster", *cmd.Cluster)
	} else if *cmd.Datacenter != "" {
		containers.AddLabelFilter(filter, "datacenter", *cmd.Datacenter)
	} else {
		// List the details of all nodes.
	}

	runningContainers, err = containers.GetDockerContainers(filter)
	if err != nil {
		panic(err)
	}

	for _, container := range runningContainers {
		tableData = append(tableData, []string{
			container.Labels[containers.DOCKERLABELCLUSTER],
			container.Names[0][1:],
			container.Labels[containers.DOCKERLABELDATACENTER],
			container.Labels[containers.DOCKERLABELIPADDRESS],
		})
	}
	printContainers(tableData)
	return nil
}

func printContainers(tableData [][]string) {
	fmt.Println()
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"CLUSTER", "HOST", "DC", "IP"})
	table.SetFooter([]string{"", "", "NODES", strconv.Itoa(len(tableData))})
	table.SetBorder(false)
	table.AppendBulk(tableData)
	table.Render()
	fmt.Println()
}
