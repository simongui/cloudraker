package commands

import (
	"fmt"
	"sync"
)

// AddBatchContext Represents context information about the AddBatch command.
type AddBatchContext struct {
	Nodes      *int
	Cluster    *string
	Datacenter *string
	HostFormat *string
	Subnet     *string
}

// NewAddBatchCommand Returns a new instance of AddBatch Command.
func NewAddBatchCommand() *Command {
	steps := []RunFunc{
		addContainersStep,
	}
	cmd := NewCommand(steps)
	return cmd
}

func addContainersStep(cmd *Command) error {
	context, _ := cmd.Context.(*AddBatchContext)
	var wg sync.WaitGroup
	wg.Add(*context.Nodes)
	cmd.Progress.Total = int64(cmd.NumberOfSteps + *context.Nodes)

	for index := 1; index <= *context.Nodes; index++ {
		host := fmt.Sprintf(*context.HostFormat, index)
		ipaddress := fmt.Sprintf("%s%d", *context.Subnet, index)

		subCommand := NewAddCommand()
		subContext := AddContext{
			Cluster:    context.Cluster,
			Datacenter: context.Datacenter,
			Host:       &host,
			IPAddress:  &ipaddress,
		}
		subCommand.Progress = cmd.Progress
		subCommand.NoProgress = cmd.NoProgress
		subCommand.IsSubCommand = true
		subCommand.SetContext(&subContext)
		cmd.Progress.Total = int64(cmd.NumberOfSteps + *context.Nodes*subCommand.NumberOfSteps)

		go func() {
			defer wg.Done()
			subCommand.RunSteps()
		}()
	}
	wg.Wait()
	return nil
}
