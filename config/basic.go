package config

type BasicArgument struct {
	Verbose bool

	*ConfigArgument
}

func NewBasicArgument() *BasicArgument {
	return &BasicArgument{
		ConfigArgument: NewConfigArgument(),
	}
}
