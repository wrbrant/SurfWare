package main

import (
	"bufio"
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

var (
	cmds               []byte
	fpath_bash_history string
	fpath_bash_aliases string
	dpath_bin          string
	commands           []string
	all                bool
	scripts            bool
	aliases            bool
	unused             bool
	add                bool
	edit               bool
	help               []string
	helpMeCmds         []string
	addMeCmds          []string
	editMeCmds         []string
)

type custom struct {
	typeOf   string //the type, whether it is an alias or a straight function
	command  string //the actual command
	category string //optional category to go under
	help     string //an optional help string to go along with it
}

var allAliases = make(map[string]custom)

func main() {
	currUser, err := user.Current()
	chk(err)
	fpath_bash_history = filepath.Join(currUser.HomeDir, ".bash_history")
	fpath_bash_aliases = filepath.Join(currUser.HomeDir, ".bash_aliases")

	dpath_bin = filepath.Join(currUser.HomeDir, "bin")
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println(helpString)
		os.Exit(0)
	}
	for i, arg := range args {
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
			help = append(help, args[i-1])
		case "add":
			add = true
		case "edit":
			edit = true
		default:
			//fairly sure this logic is screwed up
			if add {
				addMeCmds = append(addMeCmds, arg)
			} else if edit {
				editMeCmds = append(editMeCmds, arg)
			} else {
				helpMeCmds = append(helpMeCmds, arg)
			}
		}
	}
	customLoad()
	fmt.Println(allAliases["adcat"].category)
	if len(help) > 0 {
		for _, needHelp := range help {
			if _, ok := allAliases[needHelp]; !ok {
				fmt.Println(needHelp + ": command not found")
			} else {
				fmt.Printf("%s: %s\n", needHelp, allAliases[needHelp].help)
			}
		}
	}
	if add {
		customAdd(addMeCmds...)
		return
	}
	if len(editMeCmds) > 0 {
		customEdit()
		return
	}
	if all {
		scripts, aliases = true, true
	}
	if scripts {
		getScripts()
	}
	if aliases {
		getAliases()
	}
	if unused {
		getBashHistory() //these run if unused is true
		printUnused()
		return
	}

	if len(commands) > 0 {
		fmt.Println("Commands Available:")
		for _, c := range commands {
			fmt.Println("  * " + c)
		}
	}
}

func customLoad() {
	if _, err := os.Stat(fpath_bash_aliases); err != nil {
		fmt.Println(fpath_bash_aliases + " did not exist")
		os.Exit(0)
	}
	var err error
	buff, err := os.ReadFile(fpath_bash_aliases)
	chk(err)
	preg := regexp.MustCompile(`(?:#-----*(\w+)-----*|alias\s*(\w*)=((?:".*"|'.*')\s*)#(.*)\n|function\s*([\w\W]+)\(\){\s*#(.+)\n([.\n\s\w\W]+}\s*\n))`)
	matches := preg.FindAllStringSubmatch(string(buff), -1)
	cat := ""
	for _, match := range matches {
		if match[1] != "" {
			cat = match[1]
		} else if match[5] != "" {
			allAliases[match[5]] = custom{typeOf: "function", command: match[5], category: cat, help: match[6]}
		} else {
			allAliases[match[2]] = custom{typeOf: "alias", command: match[3], category: cat, help: match[4]}
		}
	}
}

func customAdd(cs ...string) {
	input := bufio.NewScanner(os.Stdin)
	readLine := func() string {
		if !input.Scan() {
			if err := input.Err(); err != nil {
				chk(err)
			}
			return ""
		}
		return strings.TrimSpace(input.Text())
	}

	if len(cs) == 0 {
		fmt.Println("Please specify a command to add")
		commandName := readLine()
		if commandName == "" {
			fmt.Println("No command specified, exiting")
			os.Exit(0)
		}
		cs = append(cs, commandName)
	}
	newAliases := make(map[string]custom)
	for _, c := range cs {
		if _, ok := allAliases[c]; ok {
			fmt.Printf("%s already exists, skipping\n", c)
		} else {
			fmt.Printf("Please enter the command for %s:\n", c)
			command := readLine()
			fmt.Printf("Please enter a help string for %s (optional):\n", c)
			help := readLine()
			fmt.Printf("Please enter a category for %s (optional):\n", c)
			category := readLine()
			newAliases[c] = custom{typeOf: "alias", command: command, category: category, help: help}
		}
	}
	fp, err := os.OpenFile(fpath_bash_aliases, os.O_APPEND|os.O_WRONLY, 0664)
	chk(err)
	defer fp.Close()
	for alias, details := range newAliases {
		if _, ok := allAliases[alias]; !ok {
			//TODO: add category support by writing the alias underneath the proper category header, and if the category doesn't exist, add a new header for it
			_, err := fmt.Fprintf(fp, "\nalias %s=\"%s\"\t#%s\n", alias, details.command, details.help)
			chk(err)
			fmt.Printf("%s added successfully\n", alias)
		}
	}
}
func customEdit() {

}

func getBashHistory() {
	if _, err := os.Stat(fpath_bash_history); err != nil {
		fmt.Println(fpath_bash_history + " did not exist")
		os.Exit(0)
	}
	var err error
	cmds, err = os.ReadFile(fpath_bash_history)
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
		if found := regexp.MustCompile(fmt.Sprintf(`\n*\s*%s\s*.*\n*`, cmd)).Find(cmds); found == nil {
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
