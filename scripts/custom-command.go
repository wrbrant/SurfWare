package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

var helpString = `
=========custom-command help=========
Helps organize custom commands -- all user made aliases and scripts on the system

Usage:
	custom-command [[-A -s -a]-u] [-h]
Flags:
	-A, --all		prints all aliases and scripts
	-s, --scripts	prints all scripts
	-a, --aliases	prints all aliases
	-u, --unused	of the commands, prints those that haven't been used in awhile. Used with -A, -s, or -a
	-h, --help		help for custom-command
	add				starts a cli for adding an alias or function
	edit			starts a cli for editing an alias or function
`

var cmds []byte
var fpath_bash_history string
var dpath_bin string
var commands []string
var all bool
var scripts bool
var aliases bool
var unused bool
var add bool
var edit bool

type custom struct {
	typ      string //alias or function
	command  string //the actual command
	category string //optional to go under
	help     string //an optional help string to go along with it
}

func main() {
	currUser, err := user.Current()
	chk(err)
	fpath_bash_history = filepath.Join(currUser.HomeDir, ".bash_history")
	dpath_bin = filepath.Join(currUser.HomeDir, "bin")

	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Printf(helpString)
	}
	for _, arg := range args {
		switch arg {
		case "-A", "--all":
			all = true
		case "-s", "--scripts":
			scripts = true
		case "-a", "--aliases":
			aliases = true
		case "-u", "--unused":
			unused = true
		case "-h", "--help":
			fmt.Println(helpString)
			os.Exit(0)
		case "add":
			add = true
		case "edit":
			edit = true
		default:
			fmt.Println("ERROR: unrecognized flag: " + arg)
			fmt.Println(helpString)
			os.Exit(1)
		}
	}
	if add {
		customAdd()
		return
	}
	if edit {
		customEdit()
		return
	}
	if all {
		getScripts()
		getAliases()
	} else {
		if scripts {
			getScripts()
		}
		if aliases {
			getAliases()
		}
	}
	if !unused {
		fmt.Println("Commands Available:")
		for _, c := range commands {
			fmt.Println("  * " + c)
		}
	} else {
		getBashHistory()
		printUnused()
	}

}

func customAdd() {

}
func customEdit() {

}

func getBashHistory() {
	if _, err := os.Stat(fpath_bash_history); err != nil {
		fmt.Println(fpath_bash_history + " did not exist")
		os.Exit(0)
	}
	var err error
	cmds, _ = os.ReadFile(fpath_bash_history)
	chk(err)
}

func getScripts() {
	items, err := ioutil.ReadDir(dpath_bin)
	chk(err)
	for _, item := range items {
		if !strings.Contains(item.Name(), ".") && !(item.Name() == "bash_aliases") && !item.IsDir() {
			commands = append(commands, item.Name())
		}
	}
}

func getAliases() {
	data, err := os.ReadFile(path.Join(dpath_bin, "bash_aliases"))
	chk(err)
	// strs := strings.Split(string(data), "\n")
	re := regexp.MustCompile(`(?:\nalias|\nfunction)\s+(\w+)`)
	matches := re.FindAllStringSubmatch(string(data), -1)
	for _, match := range matches {
		commands = append(commands, match[1])
	}
}

func printUnused() {
	var unfound []string
	for _, cmd := range commands {
		if found := regexp.MustCompile(fmt.Sprintf(`\n\s*%s\s*\n`, cmd)).Find(cmds); found == nil {
			unfound = append(unfound, cmd)
		}
	}
	if len(unfound) > 0 {
		fmt.Printf("The following scripts haven't been used in awhile:\n  * %s\n", strings.Join(unfound, "\n  * "))
	}
}

func chk(err error) {
	if err != nil {
		panic(err)
	}
}
