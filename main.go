package main

import (
	"fmt"
	"monkey/repl"
	"os"
	"os/user"
)

const VERSION = "0.0.1"

func main() {
	osUser, err := user.Current()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Monkey %s ~ %s\n", VERSION, osUser.Username)
	repl.Start(os.Stdin, os.Stdout)
}
