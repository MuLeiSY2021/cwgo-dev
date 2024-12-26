package config

import "github.com/cloudwego/cwgo/config/basic"

type BasicArgument struct {
	Verbose bool

	*basic.ConfigArgument
}

var BasicArguments *BasicArgument

func newBasicArgument() *BasicArgument {
	return &BasicArgument{
		ConfigArgument: basic.NewConfigArgument(),
	}
}

func init() {
	BasicArguments = newBasicArgument()
}
