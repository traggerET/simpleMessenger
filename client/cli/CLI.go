package cli

import (
	"bufio"
	"fmt"
	"github.com/urfave/cli"
	"simpleMessenger/client/customErrors"
	"strings"
)

type Interacter interface {
	Interact() error
}

type Interaction struct {
	app *cli.App
	r   *bufio.Reader
}

func NewInteraction(name string, r *bufio.Reader) *Interaction {
	app := cli.NewApp()
	app.Name = name
	app.CommandNotFound = func(context *cli.Context, s string) {
		fmt.Printf("Cannot recognize command: %s. Try again...", s)
	}
	return &Interaction{app, r}
}

func (i *Interaction) Interact() error {
	for {
		cmd, _ := i.r.ReadString('\n')

		cmd = strings.TrimSpace(cmd)
		args := strings.Split(cmd, " ")

		args = append([]string{i.app.Name}, args...)
		err := i.app.Run(args)

		if err == customErrors.EOI {
			return nil
		}
	}
}
