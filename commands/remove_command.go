package commands

import (
	"errors"

	"github.com/docker/engine-api/types"
	"github.com/simongui/cloudraker/containers"
)

// RemoveContext Represents context information about the Remove command.
type RemoveContext struct {
	Host      *string
	Container *types.Container
}

// NewRemoveCommand Returns a new instance of Remove Command.
func NewRemoveCommand() *Command {
	steps := []RunFunc{
		getContainerStep,
		removeContainerStep,
		//removeDockerNetwork,
		removeIPAliasStep,
		removeDNSStep,
	}
	cmd := NewCommand(steps)
	return cmd
}

func getContainerStep(cmd *Command) error {
	context, _ := cmd.Context.(*RemoveContext)
	container, err := containers.GetDockerContainer(*context.Host)
	context.Container = container

	if context.Container.ID == "" || err != nil {
		return errors.New("Unable to find container")
	}
	return nil
}

func removeContainerStep(cmd *Command) error {
	context, _ := cmd.Context.(*RemoveContext)
	err := containers.Remove(context.Container)
	return err
}

// func removeDockerNetworkStep(cmd *RemoveCommand) error {
// 	err := containers.RemoveDockerNetwork(*cmd.Host)
// 	return err
// }

func removeIPAliasStep(cmd *Command) error {
	context, _ := cmd.Context.(*RemoveContext)
	err := containers.RemoveHostsIPAlias(context.Container)
	return err
}

func removeDNSStep(cmd *Command) error {
	context, _ := cmd.Context.(*RemoveContext)
	containers.RemoveDNS(*context.Host)
	return nil
}
