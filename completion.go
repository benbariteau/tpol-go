package shex

import (
	"bufio"
	"fmt"
	"os/exec"
	"regexp"
)

func bashCompletionFuncs() map[string]string {
	cmd := exec.Command("bash", "-c", "source /usr/local/etc/bash_completion; complete -p")

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		panic(err.Error())
	}
	scanner := bufio.NewScanner(stdout)

	cmd.Start()

	completionFuncs := make(map[string]string)

	rx := regexp.MustCompile("-F ([^ ]*).* ([^ ]*)$")
	for scanner.Scan() {
		match := rx.FindStringSubmatch(scanner.Text())
		if match != nil {
			funcName := match[1]
			cmdName := match[2]
			completionFuncs[cmdName] = funcName
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err.Error())
	}

	return completionFuncs
}
