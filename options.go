package regcmd

import (
	"context"
	"fmt"
)

type options struct {
	loggerFunc     func(s string)
	contextGenFunc func() context.Context
}

type CommandOption func(o *options)

func WithContextGeneration(cgf func() context.Context) CommandOption {
	return func(o *options) {
		o.contextGenFunc = cgf
	}
}

func WithLoggerFunc(lf func(s string)) CommandOption {
	return func(o *options) {
		o.loggerFunc = lf
	}
}

var defaultOptions []CommandOption = []CommandOption{
	WithLoggerFunc(func(s string) {
		fmt.Print(s)
	}),
	WithContextGeneration(func() context.Context {
		return context.Background()
	}),
}
