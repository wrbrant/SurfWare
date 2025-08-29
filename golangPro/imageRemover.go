package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

type imagesStruct struct {
	Repo    string
	Tag     string
	ImageID string
	Created string
	Size    string
}

func listImgs() {
	cmd := exec.Command("docker", "images")
	out, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}
	if err := cmd.Start(); err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(out)
	rmImages := []string{"rmi"} //start with 'rmi' so it adds to the exec.Command cleanly
	for scanner.Scan() {
		line := scanner.Text()
		match := regexp.MustCompile(`[^\s]+\s{2}\s*[^\s]+\s{2}\s*([^\s]+)\s{2}\s*\w.+`).FindAllStringSubmatch(line, -1)
		if len(match) > 0 {
			rmImages = append(rmImages, match[0][1])
		}
	}
	fmt.Printf("The following docker images are on this machine:\n%s\nWould you like to remove all images [Y/n]", strings.Join(rmImages[1:], "\n"))
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	if input != strings.ToLower("y\n") {
		os.Exit(0)
	}
	if finalOut, err := exec.Command("docker", rmImages...).CombinedOutput(); err != nil {
		panic(err)
	} else {
		fmt.Println(string(finalOut))
	}

}

func main() {
	listImgs()
}
