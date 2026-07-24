package autocomplete

import (
	"os"
	"slices"
	"strings"
)

type CommandAutocompleter struct {
	prefix   string
	path     string
	builtins []string
	options  []string
	loaded   map[string]string
}

func (a *CommandAutocompleter) SetBuiltins(s []string) {
	a.builtins = s
}

func (a *CommandAutocompleter) SetPATH(path string) {
	a.path = path
}

func (a *CommandAutocompleter) EagerLoad() {
	commandMap := make(map[string]string, 0)

	for _, command := range a.builtins {
		_, ok := commandMap[command]
		if !ok {
			commandMap[command] = "builtin"
		}
	}

	dirs := strings.Split(a.path, string(os.PathListSeparator))
	for _, dir := range dirs {
		dirEntries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, dirEntry := range dirEntries {
			if dirEntry.IsDir() {
				continue
			}

			_, ok := commandMap[dirEntry.Name()]
			if !ok {
				commandMap[dirEntry.Name()] = dir + "/" + dirEntry.Name()
			}
		}
	}

	a.loaded = commandMap
}

func (a *CommandAutocompleter) GetLoadedCommands() map[string]string {
	return a.loaded
}

func (a *CommandAutocompleter) Match(input string) []string {
	a.prefix = input
	a.Clear()

	for commandName := range a.loaded {
		if strings.HasPrefix(commandName, a.prefix) {
			a.options = append(a.options, commandName)
		}
	}

	return a.Retrieve()
}

func (a *CommandAutocompleter) Retrieve() []string {
	slices.SortFunc(a.options, func(a, b string) int {
		return strings.Compare(a, b)
	})

	return a.options
}

func (a *CommandAutocompleter) Clear() {
	if len(a.options) > 0 {
		a.options = make([]string, 0)
	}
}

var instance *CommandAutocompleter

func GetCommandAutocompleter() *CommandAutocompleter {
	if instance == nil {
		instance = &CommandAutocompleter{}
	}

	return instance
}
