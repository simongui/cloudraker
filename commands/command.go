package commands

import (
	"fmt"
	"os"
	"strings"
	"time"

	kingpin "gopkg.in/alecthomas/kingpin.v2"

	"github.com/cheggaaa/pb"
)

// RunFunc Function that executes a command step.
type RunFunc func(cmd *Command) error

// Command Represents a command.
type Command struct {
	Progress      *pb.ProgressBar
	Results       []string
	NumberOfSteps int
	Steps         []RunFunc
	Context       interface{}
	NoProgress    *bool
	IsSubCommand  bool
}

// NewCommand Returns a new instance of Command.
func NewCommand(steps []RunFunc) *Command {
	cmd := &Command{
		// Progress:      progress,
		Steps:         steps,
		NumberOfSteps: len(steps),
	}
	return cmd
}

// SetContext Sets the commands context.
func (cmd *Command) SetContext(context interface{}) {
	cmd.Context = context
}

// Run Executes the command.
func (cmd *Command) Run(c *kingpin.ParseContext) error {
	progress := pb.New(cmd.NumberOfSteps)
	progress.ShowTimeLeft = true
	progress.ShowFinalTime = true
	progress.SetRefreshRate(time.Second * 1)
	cmd.Progress = progress

	if cmd.NoProgress == nil || *cmd.NoProgress == false {
		cmd.Progress.Start()
	}

	_, err := cmd.RunSteps()
	return err
}

// RunSteps Runs the provided steps and returns errors.
func (cmd *Command) RunSteps() (string, error) {
	var lines string

	for _, step := range cmd.Steps {
		err := step(cmd)
		if err != nil {
			if cmd.NoProgress == nil || *cmd.NoProgress == false {
				cmd.Progress.Finish()
			}
			fmt.Println(err)
			os.Exit(1)
		}
		if cmd.NoProgress == nil || *cmd.NoProgress == false {
			cmd.Progress.Increment()
		}
	}
	if cmd.NoProgress == nil || *cmd.NoProgress == false {
		if len(cmd.Results) > 0 {
			lines = strings.Join(cmd.Results[:], "")
			if cmd.IsSubCommand == false {
				cmd.Progress.FinishPrint(lines)
			}
		} else {
			cmd.Progress.Finish()
		}
	} else {
		if len(cmd.Results) > 0 {
			lines = strings.Join(cmd.Results[:], "")
			if cmd.IsSubCommand == false {
				fmt.Println(lines)
			}
		}
	}

	return "", nil
}
