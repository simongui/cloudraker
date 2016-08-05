package commands

import (
	"fmt"
	"sync"
)

// RemoveBatchContext Represents context information about the RemoveBatch command.
type RemoveBatchContext struct {
	Nodes      *int
	HostFormat *string
}

// NewRemoveBatchCommand Returns a new instance of RemoveBatch Command.
func NewRemoveBatchCommand() *Command {
	steps := []RunFunc{
		removeContainersStep,
	}
	cmd := NewCommand(steps)
	return cmd
}

func removeContainersStep(cmd *Command) error {
	context, _ := cmd.Context.(*RemoveBatchContext)
	var wg sync.WaitGroup
	wg.Add(*context.Nodes)
	cmd.Progress.Total = int64(cmd.NumberOfSteps + *context.Nodes)

	for index := 1; index <= *context.Nodes; index++ {
		host := fmt.Sprintf(*context.HostFormat, index)

		subCommand := NewRemoveCommand()
		subContext := RemoveContext{
			Host: &host,
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
