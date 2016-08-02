package commands

import (
	"errors"
	"fmt"
	"os"

	"github.com/cheggaaa/pb"
	"github.com/docker/engine-api/types"
	"github.com/simongui/cloudraker/containers"

	"gopkg.in/alecthomas/kingpin.v2"
)

type removeStepFunc func(cmd *RemoveCommand) error

// RemoveCommand Context for "remove" command
type RemoveCommand struct {
	Host      *string
	Container *types.Container
}

// Run Executes the command like action.
func (cmd *RemoveCommand) Run(c *kingpin.ParseContext) error {
	steps := []removeStepFunc{
		getContainer,
		removeContainer,
		//removeDockerNetwork,
		removeIPAlias,
	}

	progress := pb.New(len(steps)).Prefix(*cmd.Host)
	progress.ShowTimeLeft = true
	progress.ShowFinalTime = true

	pool, err := pb.StartPool(progress)
	if err != nil {
		panic(err)
	}

	for _, step := range steps {
		err := step(cmd)
		if err != nil {
			pool.Stop()
			fmt.Println(err)
			os.Exit(1)
		}
		progress.Increment()
	}
	pool.Stop()

	return nil
}

func getContainer(cmd *RemoveCommand) error {
	var err error
	cmd.Container, err = containers.GetDockerContainer(*cmd.Host)
	if cmd.Container.ID == "" || err != nil {
		return errors.New("Unable to find container")
	}
	return nil
}

func removeContainer(cmd *RemoveCommand) error {
	err := containers.Remove(cmd.Container)
	return err
}

// func removeDockerNetwork(cmd *RemoveCommand) error {
// 	err := containers.RemoveDockerNetwork(*cmd.Host)
// 	return err
// }

func removeIPAlias(cmd *RemoveCommand) error {
	err := containers.RemoveHostsIPAlias(cmd.Container)
	return err
}
