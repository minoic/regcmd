package regcmd

import (
	"context"
	"fmt"
)

type options struct {
	loggerFunc     func(s string)
	contextGenFunc func() context.Context
	pool           chan struct{}
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

func WithPoolSize(psize uint) CommandOption {
	return func(o *options) {
		o.pool = make(chan struct{}, psize)
	}
}

var defaultOptions []CommandOption = []CommandOption{
	WithLoggerFunc(func(s string) {
		fmt.Println("regcmd: ", s)
	}),
	WithContextGeneration(func() context.Context {
		return context.Background()
	}),
	WithPoolSize(1),
}
