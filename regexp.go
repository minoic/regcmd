package regcmd

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"regexp"
)

// register a command
// re eg. `command (.*) flag (.*)`
// handler: args:the params matched by re
func Register(re string, handler func(args []string) error) {
	instance.register(re, handler)
}

// listen io.Reader and will block the routine
// Normally it can be os.Stdin
func Listen(stream io.Reader) {
	instance.listen(stream)
}

type command struct {
	Re      *regexp.Regexp
	Handler func(args []string) error
}

type manager struct {
	C []command
}

var instance manager

func (this *manager) register(re string, handler func(args []string) error) {
	rec, err := regexp.Compile(re)
	if err != nil {
		panic(err)
	}
	c := command{
		Re:      rec,
		Handler: handler,
	}
	if this == nil {
		this.C = append([]command{}, c)
	} else {
		this.C = append(this.C, c)
	}
}

func (this *manager) handle(input string) error {
	for _, c := range this.C {
		if args := c.Re.FindStringSubmatch(input); len(args)-1 == c.Re.NumSubexp() {
			return c.Handler(args[1:])
		}
	}
	return errors.New("invalid command: " + input)
}

func (this *manager) listen(stream io.Reader) {
	reader := bufio.NewReader(stream)
	for {
		var input string
		b, _, _ := reader.ReadLine()
		input = string(b)
		if len(input) == 0 {
			continue
		}
		err := this.handle(input)
		if err != nil {
			fmt.Println(err)
		}
	}
}
