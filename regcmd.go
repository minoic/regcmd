package regcmd

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"
)

// Register register a command
//
// re: eg. `command (.*) flag (.*)`
//
// names: if you have n arguments in re, you should have n names at least.
// In additional, if you have n+1 names ,the last one will be the introduction of this command
//
// handler: args:the params matched by re
func Register(re string, names []string, handler ...Handler) error {
	return instance.register(re, names, handler...)
}

// ShouldRegister registers a command and will panic while an error received
func ShouldRegister(re string, names []string, handler ...Handler) {
	err := instance.register(re, names, handler...)
	if err != nil {
		panic(err)
	}
}

// Listen listen io.Reader and will block the routine
// Normally it can be os.Stdin
//
// example: go regcmd.Listen(os.Stdin)
func Listen(stream io.Reader, optfuncs ...CommandOption) {
	instance.listen(stream, optfuncs...)
}

type command struct {
	re       *regexp.Regexp
	words    int
	desc     string
	intro    string
	handlers []Handler
}

type manager struct {
	rwLock   sync.RWMutex
	commands []command
	opts     options
	optfuncs []CommandOption
}

type Handler func(ctx *Context, args []string)

type Context struct {
	c       context.Context
	aborted bool
}

func (c Context) Deadline() (deadline time.Time, ok bool) {
	return c.c.Deadline()
}

func (c Context) Done() <-chan struct{} {
	return c.c.Done()
}

func (c Context) Err() error {
	return c.c.Err()
}

func (c Context) Value(key interface{}) interface{} {
	return c.c.Value(key)
}

func (c Context) Set(key, value interface{}) {
	c.c = context.WithValue(c.c, key, value)
}

func (c Context) Abort() {
	c.aborted = true
}

var (
	instance manager
	helper   = make(map[string][]*command)
	once     sync.Once
)

func (this *manager) register(re string, names []string, handler ...Handler) error {
	rec, err := regexp.Compile(re)
	if err != nil {
		return errors.New(re + err.Error())
	}
	if len(names) < rec.NumSubexp() {
		return errors.New(re + " not enough names for parenthesized subexpressions in this Regexp")
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
		re:       rec,
		words:    len(splts),
		desc:     buf.String(),
		handlers: handler,
	}
	if count < len(names) {
		c.intro = names[count]
	}
	helper[splts[0]] = append(helper[splts[0]], &c)
	if splts[0] != "help" && len(helper[splts[0]]) == 1 {
		err = instance.register(splts[0]+" help", []string{"To get this help"}, func(ctx *Context, args []string) {
			this.opts.loggerFunc(fmt.Sprintf("---- %s help ----\n", splts[0]))
			var buf bytes.Buffer
			for _, c := range helper[splts[0]] {
				buf.WriteString(c.desc)
				if len(c.intro) != 0 {
					for i := 1; i <= len(helper[splts[0]][0].desc)-len(c.desc); i++ {
						buf.WriteByte(' ')
					}
					buf.WriteString("// ")
					buf.WriteString(c.intro)
				}
				this.opts.loggerFunc(buf.String())
				buf.Reset()
			}
		})
		if err != nil {
			return err
		}
		once.Do(func() {
			err = instance.register("help", []string{"To get all commands help"}, func(ctx *Context, args []string) {
				this.opts.loggerFunc("---- all commands help ----")
				for k, _ := range helper {
					var buf bytes.Buffer
					for _, c := range helper[k] {
						buf.WriteString(c.desc)
						if len(c.intro) != 0 {
							for i := 1; i <= len(helper[k][0].desc)-len(c.desc); i++ {
								buf.WriteByte(' ')
							}
							buf.WriteString("// ")
							buf.WriteString(c.intro)
						}
						this.opts.loggerFunc(buf.String())
						buf.Reset()
					}
				}
			})
			if err != nil {
				panic(err)
			}
		})
	}
	this.rwLock.Lock()
	if this == nil {
		this.commands = append([]command{}, c)
	} else {
		this.commands = append(this.commands, c)
	}
	this.rwLock.Unlock()
	return nil
}

func (this *manager) handle(input string) string {
	splts := strings.Split(input, " ")
	for _, c := range this.commands {
		if c.words != len(splts) {
			continue
		}
		ctx := &Context{
			c:       this.opts.contextGenFunc(),
			aborted: false,
		}
		if args := c.re.FindStringSubmatch(input); len(args)-1 == c.re.NumSubexp() {
			for i := range c.handlers {
				c.handlers[i](ctx, args[1:])
				if ctx.aborted {
					return ""
				}
			}
			return ""
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
	return "Invalid command: " + input + ` **Type <help> for commands help`
}

func (this *manager) listen(stream io.Reader, optfuncs ...CommandOption) {
	for i := range defaultOptions {
		defaultOptions[i](&this.opts)
	}
	for i := range optfuncs {
		optfuncs[i](&this.opts)
	}
	for k := range helper {
		sort.Slice(helper[k], func(i, j int) bool {
			return len(helper[k][i].desc) > len(helper[k][j].desc)
		})
	}
	reader := bufio.NewReader(stream)
	for {
		var input string
		b, _, _ := reader.ReadLine()
		input = string(b)
		if len(input) == 0 {
			continue
		}
		this.opts.pool <- struct{}{}
		go func() {
			defer func() {
				<-this.opts.pool
			}()
			if ret := this.handle(input); len(ret) != 0 {
				this.opts.loggerFunc(ret)
			}
		}()
	}
}
