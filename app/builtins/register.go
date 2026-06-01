package builtins

type Builtin func(args []string) error

var Builtins = map[string]Builtin{}

func Register(name string, command Builtin) {
	_, exists := Builtins[name]

	if !exists {
		// silently ignore the command registered twice for now.
		Builtins[name] = command
	}
}
