package main

import (
	"os"

	_ "github.com/go-sql-driver/mysql" // Required to use the MySQL driver.
	"github.com/simongui/cloudraker/commands"
	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	app := kingpin.New("cluster", "Cluster management application.")

	lsCommand := &commands.LsCommand{}
	ls := app.Command("ls", "List nodes in a specific or all clusters.").Action(lsCommand.Run)
	lsCommand.Cluster = ls.Flag("cluster", "Cluster to list.").String()
	lsCommand.Datacenter = ls.Flag("datacenter", "Datacenter to list.").String()
	lsCommand.Host = ls.Flag("host", "Host to list.").String()

	addCommand := &commands.AddCommand{}
	add := app.Command("add", "Add a node to the specified cluster and datacenter.").Action(addCommand.Run)
	addCommand.Cluster = add.Flag("cluster", "Cluster to add the new node to.").String()
	addCommand.Datacenter = add.Flag("datacenter", "Datacenter to add the new node to.").String()
	addCommand.Host = add.Flag("host", "Host to associate with the new node.").String()
	addCommand.IPAddress = add.Flag("ipaddress", "IP address to associate with the new node.").String()

	removeCommand := &commands.RemoveCommand{}
	remove := app.Command("remove", "Remove a node from the cluster.").Action(removeCommand.Run)
	removeCommand.Host = remove.Arg("host", "Host to remove from the cluster.").String()

	kingpin.MustParse(app.Parse(os.Args[1:]))
}
