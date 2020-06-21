package regcmd

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"regexp"
	"sort"
	"strings"
)

// register a command
// re eg. `command (.*) flag (.*)`
// handler: args:the params matched by re
func Register(re string, names []string, handler func(args []string) string) {
	instance.register(re, names, handler)
}

// listen io.Reader and will block the routine
// Normally it can be os.Stdin
func Listen(stream io.Reader) {
	instance.listen(stream)
}

type command struct {
	Re      *regexp.Regexp
	Desc    string
	Intro   string
	Handler func(args []string) string
}

type manager struct {
	C []command
}

var (
	instance manager
	helper   = make(map[string][]*command)
)

func (this *manager) register(re string, names []string, handler func(args []string) string) {
	rec, err := regexp.Compile(re)
	if err != nil {
		panic(re + err.Error())
	}
	if len(names) < rec.NumSubexp() {
		panic(re + " not enough names for parenthesized subexpressions in this Regexp")
	}
	splts := strings.Split(re, " ")
	var buf bytes.Buffer
	count := 0
	for _, v := range splts {
		// only supports "(.*)"
		if v == "(.*)" {
			buf.WriteByte('<')
			buf.WriteString(names[count])
			buf.WriteString("> ")
			count++
		} else {
			buf.WriteString(v)
			buf.WriteByte(' ')
		}
	}
	c := command{
		Re:      rec,
		Desc:    buf.String(),
		Handler: handler,
	}
	if count < len(names) {
		c.Intro = names[count]
	}
	helper[splts[0]] = append(helper[splts[0]], &c)
	sort.Slice(helper[splts[0]], func(i, j int) bool {
		return len(helper[splts[0]][i].Desc) > len(helper[splts[0]][j].Desc)
	})
	if len(helper[splts[0]]) == 1 {
		instance.register(splts[0]+" help", []string{"To get this help"}, func(args []string) string {
			fmt.Printf("---- %s help ----\n", splts[0])
			var buf bytes.Buffer
			for _, c := range helper[splts[0]] {
				buf.WriteString(c.Desc)
				if len(c.Intro) != 0 {
					for i := 1; i <= len(helper[splts[0]][0].Desc)-len(c.Desc); i++ {
						buf.WriteByte(' ')
					}
					buf.WriteString("// ")
					buf.WriteString(c.Intro)
				}
				fmt.Println(buf.String())
				buf.Reset()
			}
			return ""
		})
	}
	if this == nil {
		this.C = append([]command{}, c)
	} else {
		this.C = append(this.C, c)
	}
}

func (this *manager) handle(input string) string {
	for _, c := range this.C {
		if args := c.Re.FindStringSubmatch(input); len(args)-1 == c.Re.NumSubexp() {
			return c.Handler(args[1:])
		}
	}
	index := strings.IndexByte(input, ' ')
	var s string
	if index == -1 {
		s = input
	} else {
		s = input[:index]
	}
	if _, ok := helper[s]; ok {
		return fmt.Sprintf("Type <%s help> for more help", s)
	}
	return "Invalid command: " + input
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
		if ret := this.handle(input); len(ret) != 0 {
			fmt.Println(ret)
		}
	}
}
