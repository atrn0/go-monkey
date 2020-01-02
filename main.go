package main

import (
	"fmt"
	"github.com/atrn0/go-monkey/repl"
	"os"
	"os/user"
)

func main() {
	currentUser, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello %s!! This is Monkey REPL!\n", currentUser.Username)
	repl.Start(os.Stdin, os.Stdout)
}
