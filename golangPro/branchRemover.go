package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func PANIC(err error) {
	if err != nil {
		panic(err)
	}
}

// GetGitBranch returns the output of "git branch"
func GetGitBranch() {
	out, err := exec.Command("git", "fetch", "-p").CombinedOutput()
	PANIC(err)
	fmt.Println(string(out))
	local, err := exec.Command("git", "branch").CombinedOutput()
	PANIC(err)
	llines := strings.Split(string(local), "\n")

	remote, err := exec.Command("git", "branch", "--remotes").CombinedOutput()
	PANIC(err)
	rlines := strings.Split(string(remote), "\n")

	fmt.Print("\033[94mlocal branches:\033[0m\n", strings.Join(llines, "\n"), "\n")
	fmt.Print("\033[94mremote branches:\n\033[91m", strings.Join(rlines, "\n"), "\033[0m\n")
	removals := []string{}
	for _, lline := range llines {
		lline = strings.Trim(strings.ReplaceAll(lline, "* ", ""), " ")
		for j, rline := range rlines {
			if strings.Contains(rline, "/"+lline) {
				break
			}
			if j == len(rlines)-1 {
				removals = append(removals, lline)
			}
		}
	}
	if len(removals) == 0 {
		fmt.Println("\033[94mThere is nothing to remove. Every branch has a corresponding remote.\033[0m")
		os.Exit(0)
	}
	fmt.Print("\033[94mThe following lines are set for removal:\033[0m\n  ", strings.Join(removals, "\n  "), "\n\n")
	fmt.Println("\033[93mConfirm removal? This is irreversible. [Y/n]\033[0m")
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		panic("scanner broke?")
	}
	if strings.ToUpper(scanner.Text()) == "Y" {
		args := append([]string{"branch", "-D"}, removals...)
		out, err := exec.Command("git", args...).CombinedOutput()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(out))
	} else {
		fmt.Println("exiting.")
		os.Exit(0)
	}
}

func main() {
	GetGitBranch()
}
