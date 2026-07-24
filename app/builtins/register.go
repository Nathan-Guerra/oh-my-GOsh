package builtins

type Response struct {
	Out        string
	Err        string
	ShouldExit bool
	ExitSignal int
}

type Builtin interface {
	Exec(args []string) *Response
}

var Builtins = map[string]Builtin{}

func Register(name string, command Builtin) {
	_, exists := Builtins[name]

	if !exists {
		// silently ignore the command registered twice for now.
		Builtins[name] = command
	}
}
