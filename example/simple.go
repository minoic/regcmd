package main

import (
	"fmt"
	"github.com/MinoIC/regcmd"
	"os"
)

func main() {
	_ = regcmd.Register("show", []string{"say hello world"}, func(args []string) {
		fmt.Println("hello world")
	})
	_ = regcmd.Register("show (.*)", []string{"user", "say hello to the given user name"}, func(args []string) {
		fmt.Println("hello", args[0])
	})
	regcmd.Listen(os.Stdin)
}
