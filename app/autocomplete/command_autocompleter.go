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

func (a *CommandAutocompleter) LargestCommonPrefix(input string) []byte {
	t := strings.Split(strings.TrimLeft(input, " "), " ")
	var pos int
	if len(t) == 1 { // search command
		pos = len(input)
	} else {
		pos = len(t[len(t)-1])
	}

	// search files
	valid := true
	prefix := make([]byte, 0)
	var target byte
	for valid {
		if pos >= len(a.options[0]) {
			break
		}
		target = a.options[0][pos:][0]
		for _, name := range a.options {
			if pos >= len(name) || name[pos] != target {
				valid = false
				break
			}
		}
		if !valid {
			break
		}
		prefix = append(prefix, target)
		pos++
	}

	return prefix
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

			info, err := dirEntry.Info()
			if err != nil {
				continue
			}

			_, ok := commandMap[dirEntry.Name()]
			if info.Mode().Perm()&0111 != 0 && !ok {
				commandMap[dirEntry.Name()] = dir + "/" + dirEntry.Name()
			}
		}
	}

	a.loaded = commandMap
}

func (a *CommandAutocompleter) GetLoadedCommands() map[string]string {
	return a.loaded
}

func (a *CommandAutocompleter) Match(input string) ([]string, int) {
	t := strings.Split(strings.TrimLeft(input, " "), " ")
	a.prefix = input
	a.Clear()

	if len(t) == 1 { // search command
		for commandName := range a.loaded {
			if strings.HasPrefix(commandName, a.prefix) {
				a.options = append(a.options, commandName)
			}
		}

		return a.Retrieve(), len(a.prefix)
	}

	// search files
	target := t[len(t)-1]

	dir, err := os.Getwd()
	if err != nil {
		return a.Retrieve(), 0
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return a.Retrieve(), 0
	}

	for _, entry := range entries {
		if entry.IsDir() == false && strings.HasPrefix(entry.Name(), target) {
			a.options = append(a.options, entry.Name())
		}
	}

	return a.Retrieve(), len(target)
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
