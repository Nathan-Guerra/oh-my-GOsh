package autocomplete

import (
	"fmt"
	"os"
	"strings"
)

type Autocompleter struct {
	PATH                    string
	prefix                  string
	builtins                []string
	computedPATHExecutables []string
	options                 []string
}

func (a *Autocompleter) SetBuiltins(s []string) {
	a.builtins = s
}

func (a *Autocompleter) SetPATH(path string) {
	fmt.Printf("\r\nset path with len: %d\r\n", len(path))
	a.PATH = path
}

func (a *Autocompleter) EagerLoadPathCommands() {
	dirs := strings.Split(a.PATH, string(os.PathListSeparator))
	for _, dir := range dirs {
		dirEntries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, dirEntry := range dirEntries {
			if dirEntry.IsDir() {
				continue
			}
			a.computedPATHExecutables = append(a.computedPATHExecutables, dirEntry.Name())
		}
	}
}

func (a *Autocompleter) MatchBuiltins() {
	for _, command := range a.builtins {
		if strings.HasPrefix(command, a.prefix) {
			a.options = append(a.options, command)
		}
	}
}

func (a *Autocompleter) MatchPATH() {
	for _, command := range a.computedPATHExecutables {
		if strings.HasPrefix(command, a.prefix) {
			a.options = append(a.options, command)
		}
	}
}

func (a *Autocompleter) Match(input string) []string {
	a.prefix = input
	a.Clear()
	a.MatchBuiltins()
	a.MatchPATH()
	return a.options
}

func (a *Autocompleter) Retrieve() []string {
	return a.options
}

func (a *Autocompleter) Clear() {
	if len(a.options) > 0 {
		a.options = make([]string, 0)
	}
}
