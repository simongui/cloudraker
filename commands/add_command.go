package commands

import (
	"fmt"
	"net"

	"github.com/simongui/cloudraker/containers"
	"github.com/simongui/cloudraker/mysql"
)

// AddContext Represents context information about the Add command.
type AddContext struct {
	// Progress   *pb.ProgressBar
	Cluster    *string
	Datacenter *string
	Host       *string
	IPAddress  *string
	host       string
	mysqlPort  int
	sshPort    int
	serverID   string
	readOnly   bool
}

// NewAddCommand Returns a new instance of Add Command.
func NewAddCommand() *Command {
	steps := []RunFunc{
		containerExistsStep,
		addContainerStep,
		addToDockerNetworkStep,
		addIPAliasStep,
		addDNSStep,
		startSSHStep,
		setReplicationGrants,
		getMySQLResults,
	}

	cmd := NewCommand(steps)
	return cmd
}

func containerExistsStep(cmd *Command) error {
	context, _ := cmd.Context.(*AddContext)

	// Check if the node is already added.
	exists := containers.Exists(*context.Host)
	if exists {
		return fmt.Errorf("%s is already a member of an existing cluster", *context.Host)
	}
	return nil
}

func addContainerStep(cmd *Command) error {
	context, _ := cmd.Context.(*AddContext)
	ip := net.ParseIP(*context.IPAddress)
	ip = ip.To4()
	context.mysqlPort = 33300 + int(ip[3])
	context.sshPort = 22200 + int(ip[3])

	err := containers.Add(*context.Cluster, *context.Datacenter, *context.Host, *context.IPAddress, context.mysqlPort, context.sshPort)
	return err
}

func addToDockerNetworkStep(cmd *Command) error {
	context, _ := cmd.Context.(*AddContext)

	containers.AddToDockerNetwork(*context.Cluster, *context.Host)
	return nil
}

func addIPAliasStep(cmd *Command) error {
	context, _ := cmd.Context.(*AddContext)

	containers.AddIPAlias(*context.IPAddress)
	return nil
}

func addDNSStep(cmd *Command) error {
	context, _ := cmd.Context.(*AddContext)
	containers.AddDNS(*context.Host, *context.IPAddress)
	return nil
}

func startSSHStep(cmd *Command) error {
	context, _ := cmd.Context.(*AddContext)
	containers.StartSSH(*context.Host)
	return nil
}

func setReplicationGrants(cmd *Command) error {
	context, _ := cmd.Context.(*AddContext)
	err := mysql.SetReplicationGrants(*context.Host, context.mysqlPort, "root", "password", "30s", "repl_user", "password")
	if err != nil {
		return err
	}
	return nil
}

func getMySQLResults(cmd *Command) error {
	context, _ := cmd.Context.(*AddContext)

	var err error
	context.host, err = mysql.GetHostname(*context.Host, context.mysqlPort, "root", "password", "30s")
	if err != nil {
		return err
	}
	context.serverID, err = mysql.GetServerID(*context.Host, context.mysqlPort, "root", "password", "30s")
	if err != nil {
		return err
	}

	context.readOnly, err = mysql.GetReadOnly(*context.Host, context.mysqlPort, "root", "password", "30s")
	if err != nil {
		return err
	}

	finishText := fmt.Sprintf("MySQL running \n\thost: %s\n\tid: %s\n\tread_only: %t\n",
		context.host,
		context.serverID,
		context.readOnly)

	cmd.Results = append(cmd.Results, finishText)
	return nil
}
