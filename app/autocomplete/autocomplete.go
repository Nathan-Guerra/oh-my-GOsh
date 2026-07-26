package autocomplete

import (
	"io/fs"
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

func (a *CommandAutocompleter) MatchCommand() ([]string, int) {
	for commandName := range a.loaded {
		if strings.HasPrefix(commandName, a.prefix) {
			a.options = append(a.options, commandName+" ")
		}
	}
	return a.Retrieve(), len(a.prefix)
}

func (a *CommandAutocompleter) MatchFile() ([]string, int) {
	// fmt.Printf("\r\na.prefix: %s\r\n", a.prefix)
	var curDir string
	if strings.HasPrefix(a.prefix, string(os.PathSeparator)) {
		curDir = string(os.PathSeparator)
		a.prefix = a.prefix[1:]
	} else {
		curDir = os.Getenv("PWD")
	}

	// fmt.Printf("\r\ncurrent directory: %s\r\n", curDir)

	var target string
	// divido em possiveis diretorios
	directories := strings.Split(a.prefix, string(os.PathSeparator))
	// ultimo index é o alvo do autocomplete
	if len(directories) > 1 {
		target = directories[len(directories)-1]
	} else {
		target = directories[0]
	}
	directories = directories[:len(directories)-1]

	// fmt.Printf("\r\ndirectories: %v | target: %s\r\n", directories, target)
	if len(directories) == 0 {
		entries, err := os.ReadDir(curDir)
		if err != nil {
			return make([]string, 0), 0
		}
		// fmt.Printf("\r\nsearching entries: %v", entries)
		// search only entries
		for _, entry := range entries {
			if strings.HasPrefix(entry.Name(), target) {
				name := entry.Name()
				if entry.IsDir() {
					name += string(os.PathSeparator)
				} else {
					name += " "
				}
				a.options = append(a.options, name)
			}
		}
	} else {
		root := strings.Join(directories, string(os.PathSeparator))
		currentDirectory := os.DirFS(curDir)

		// fmt.Printf("\r\npath: %s", root)

		fs.WalkDir(currentDirectory, root, func(path string, d fs.DirEntry, err error) error {
			if path == root {
				return nil
			}

			// fmt.Printf("\r\npath: %s | d: %s", path, d.Name())
			if err != nil {
				// fmt.Printf("\r\nerror: %s\r\n", err.Error())
				return err
			}

			if strings.HasPrefix(d.Name(), target) {
				name := d.Name()
				if d.IsDir() {
					name += string(os.PathSeparator)
				} else {
					name += " "
				}
				a.options = append(a.options, name)
			}

			// fmt.Printf("\r\npath: %s | d: %s", path, d.Name())

			if d.IsDir() {
				// fmt.Printf("\r\nskip dir for: %s | %s\r\n", path, d.Name())
				return fs.SkipDir
			}
			return nil
		})
	}

	return a.Retrieve(), len(target)
}

// the content of a.options is not exactly the match instead it fills with the
// content that will be put at the terminal if this match is accepted.
//
// IMPORTANT: The match might contain a trailing space or a path separator
// depending on what it is matching against.
func (a *CommandAutocompleter) Match(input string) ([]string, int) {
	t := strings.Split(strings.TrimLeft(input, " "), " ")
	a.prefix = input
	a.Clear()
	if len(t) == 1 {
		return a.MatchCommand()
	}

	a.prefix = t[len(t)-1]
	return a.MatchFile()
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
