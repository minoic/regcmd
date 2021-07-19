package main

import (
	"context"
	"fmt"
	"github.com/MinoIC/regcmd"
	"math/rand"
	"os"
	"time"
)

func main() {
	regcmd.ShouldRegister("show", []string{"say hello world"}, func(ctx *regcmd.Context, args []string) {
		fmt.Println("hello world")
	})
	regcmd.ShouldRegister("show (.*)", []string{"user", "say hello to the given user name"}, func(ctx *regcmd.Context, args []string) {
		fmt.Println("hello", args[0])
	})
	err := regcmd.Register("sleep", []string{}, func(ctx *regcmd.Context, args []string) {
		id := rand.Int()
		fmt.Println(id, "starts to sleep")
		time.Sleep(5 * time.Second)
		fmt.Println(id, "sleeped 5 seconds")
	})
	if err != nil {
		fmt.Println(err)
	}
	regcmd.Listen(os.Stdin, regcmd.WithPoolSize(2),
		regcmd.WithLoggerFunc(func(s string) {
			fmt.Println("log: ", s)
		}),
		regcmd.WithContextGeneration(func() context.Context {
			return context.Background()
		}),
	)
}
