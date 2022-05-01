package cli

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/urfave/cli"
	"math/rand"
	"simpleMessenger/client/customErrors"
	"simpleMessenger/client/handlers"
	"time"
)

const (
	uIDRange = 100000
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func NewMenuInteraction(name string, r *bufio.Reader) *Interaction {
	i := NewInteraction(name, r)
	i.app.Commands = []cli.Command{
		{
			Name:    "exit",
			Aliases: []string{"e"},
			Usage:   "use it to exit app",
			Action: func(c *cli.Context) error {
				fmt.Println("See you soon...")
				return &customErrors.EndOfInteraction{}
			},
		},
		{
			Name:    "login",
			Aliases: []string{"l"},
			Usage:   "login as new user to server with specified ip",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "ip",
					Value: "localhost",
					Usage: "server ip to connect to",
				},
				&cli.StringFlag{
					Name:  "user",
					Value: fmt.Sprintf("defaultUser%d", rand.Intn(uIDRange)),
					Usage: "username in chats",
				},
			},
			Action: func(c *cli.Context) error {
				userName := c.String("user")
				host := c.String("ip")
				usr, err := handlers.Login(userName, host)
				if err != nil {
					return err
				}

				defer func() {
					usr.C.Close()
				}()
				defer func() {
					usr.C.Call("Msgr.LeaveRoom", *usr.Status, usr.Status)
				}()

				userMenu := NewUserInteraction("authorized", usr, r)
				err = userMenu.Interact()
				errors.Unwrap(err)
				return err
			},
		},
	}
	return i
}
