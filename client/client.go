package main

import (
	"bufio"
	"os"
	"simpleMessenger/client/cli"
)

const (
	WelcomeMsg = "Welcome to my simpleMessenger app.\n"
)

func main() {
	r := bufio.NewReader(os.Stdin)

	menu := cli.NewMenuInteraction("menu", r)

	menu.Interact()

}
