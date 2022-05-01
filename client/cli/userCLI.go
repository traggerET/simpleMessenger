package cli

import (
	"bufio"
	"fmt"
	"github.com/urfave/cli"
	"simpleMessenger/client/customErrors"
	"simpleMessenger/client/handlers"
	"simpleMessenger/client/user"
)

func NewUserInteraction(name string, u *user.User, r *bufio.Reader) Interacter {
	i := NewInteraction(name, r)
	i.app.Commands = []cli.Command{
		{
			Name:    "exit",
			Aliases: []string{"e"},
			Usage:   "use it to sign out",
			Action: func(c *cli.Context) error {
				fmt.Println("See you soon...")

				return &customErrors.EndOfInteraction{}
			},
		},
		{
			Name:    "create",
			Aliases: []string{"c"},
			Usage:   "create new chat",
			Action: func(c *cli.Context) error {
				err := handlers.CreateChat(u)
				if err != nil {
					return err
				}
				chatMenu := NewChatInteraction("authorized", u, r)
				err = chatMenu.Interact()
				return err
			},
		},
		{
			Name:    "join",
			Aliases: []string{"j"},
			Usage:   "join existing chat",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "chatId",
					Usage: "chat you want join to",
				},
			},
			Action: func(c *cli.Context) error {
				chat := c.String("chatId")
				err := handlers.JoinChat(u, chat)
				if err != nil {
					return err
				}
				chatMenu := NewChatInteraction("authorized", u, r)
				err = chatMenu.Interact()
				return err
			},
		},
	}
	return i
}
