package main

import (
	"os"

	_ "github.com/go-sql-driver/mysql" // Required to use the MySQL driver.
	"github.com/simongui/cloudraker/commands"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	removeHost *string
)

func main() {
	app := kingpin.New("cluster", "Cluster management application.")

	lsCommand := &commands.LsCommand{}
	ls := app.Command("ls", "List nodes in a specific or all clusters.").Action(lsCommand.Run)
	lsCommand.Cluster = ls.Flag("cluster", "Cluster to list.").String()
	lsCommand.Datacenter = ls.Flag("datacenter", "Datacenter to list.").String()
	lsCommand.Host = ls.Flag("host", "Host to list.").String()

	addCommand := commands.NewAddCommand()
	add := app.Command("add", "Add a node to the specified cluster and datacenter.").Action(addCommand.Run)
	addCommand.NoProgress = add.Flag("noprogress", "Suppress and hide the progress bar.").Bool()
	addContext := commands.AddContext{
		Cluster:    add.Flag("cluster", "Cluster to add the new node to.").String(),
		Datacenter: add.Flag("datacenter", "Datacenter to add the new node to.").String(),
		Host:       add.Flag("host", "Host to associate with the new node.").String(),
		IPAddress:  add.Flag("ipaddress", "IP address to associate with the new node.").String(),
	}
	addCommand.SetContext(&addContext)

	addBatchCommand := commands.NewAddBatchCommand()
	addBatch := app.Command("addbatch", "Add a batch of nodes to the specified cluster and datacenter.").Action(addBatchCommand.Run)
	addBatchCommand.NoProgress = addBatch.Flag("noprogress", "Suppress and hide the progress bar.").Bool()
	addBatchContext := commands.AddBatchContext{
		Nodes:      addBatch.Flag("nodes", "Number of nodes to add the cluster and datacenter.").Int(),
		Cluster:    addBatch.Flag("cluster", "Cluster to add the new node to.").String(),
		Datacenter: addBatch.Flag("datacenter", "Datacenter to add the new node to.").String(),
		HostFormat: addBatch.Flag("hostformat", "Host format to associate with the new nodes.").String(),
		Subnet:     addBatch.Flag("subnet", "Subnet to associate with the new nodes.").String(),
	}
	addBatchCommand.SetContext(&addBatchContext)

	removeCommand := commands.NewRemoveCommand()
	remove := app.Command("remove", "Remove a node from the cluster.").Action(removeCommand.Run)
	removeCommand.NoProgress = remove.Flag("noprogress", "Suppress and hide the progress bar.").Bool()
	removeContext := commands.RemoveContext{Host: remove.Arg("host", "Host to remove from the cluster.").String()}
	removeCommand.SetContext(&removeContext)

	removeBatchCommand := commands.NewRemoveBatchCommand()
	removeBatch := app.Command("removebatch", "Add a batch of nodes to the specified cluster and datacenter.").Action(removeBatchCommand.Run)
	removeBatchCommand.NoProgress = removeBatch.Flag("noprogress", "Suppress and hide the progress bar.").Bool()
	removeBatchContext := commands.RemoveBatchContext{
		Nodes:      removeBatch.Flag("nodes", "Number of nodes to add the cluster and datacenter.").Int(),
		HostFormat: removeBatch.Flag("hostformat", "Host format to associate with the new nodes.").String(),
	}
	removeBatchCommand.SetContext(&removeBatchContext)

	kingpin.MustParse(app.Parse(os.Args[1:]))
}
