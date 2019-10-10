package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/rsb/monkey_interpreter/repl"
)

func main() {
	user, err := user.Current()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
	}

	fmt.Printf("Hello %s! This is the Monkey programming language!\n", user.Username)
	fmt.Printf("Feel free to type in commands\n")
	repl.Start(os.Stdin, os.Stdout)
}
