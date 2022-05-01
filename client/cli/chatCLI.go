package cli

import (
	"bufio"
	"fmt"
	"github.com/urfave/cli"
	"simpleMessenger/client/customErrors"
	"simpleMessenger/client/handlers"
	"simpleMessenger/client/user"
)

type ChatInteraction struct {
	Interacter
	u *user.User
}

func (i *ChatInteraction) Interact() error {
	fmt.Printf("chat's id: %s\n", i.u.Status.ChatId)
	err := i.Interacter.Interact()
	defer handlers.LeaveChat(i.u)
	return err
}

func NewChatInteraction(name string, u *user.User, r *bufio.Reader) Interacter {
	ni := NewInteraction(name, r)
	ni.app.Commands = []cli.Command{
		{
			Name:    "exit",
			Aliases: []string{"e"},
			Usage:   "use it to exit chat",
			Action: func(c *cli.Context) error {
				fmt.Println("See you soon...")
				return &customErrors.EndOfInteraction{}
			},
		},
		{
			Name:    "post",
			Aliases: []string{"p"},
			Usage:   "post some message in chat",
			Action: func(c *cli.Context) error {
				msg, _ := r.ReadString('\n')

				if msg != "" {
					err := handlers.PostMsg(u, msg)
					if err != nil {
						fmt.Println(err.Error())
					}
				}

				return nil
			},
		},
		{
			Name:    "update",
			Aliases: []string{"u"},
			Usage:   "get unread message history",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "file",
					Value: "hist.txt",
					Usage: "file to save hist to",
				},
			},
			Action: func(c *cli.Context) error {
				hist, err := handlers.GetUnread(u)
				if err != nil {
					return err
				}
				handlers.SaveHist(hist, c.String("file"), true)
				return nil
			},
		},
		{
			Name:    "hist",
			Aliases: []string{"hist"},
			Usage:   "get full message history",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "file",
					Value: "hist.txt",
					Usage: "file to save hist to",
				},
			},
			Action: func(c *cli.Context) error {
				hist, err := handlers.GetHistory(u)
				if err != nil {
					return err
				}
				handlers.SaveHist(hist, c.String("file"), false)
				return nil
			},
		},
	}
	i := ChatInteraction{ni, u}
	return &i
}
