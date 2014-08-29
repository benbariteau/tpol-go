package main

import (
	"fmt"
	"github.com/GeertJohan/go.linenoise"
	"github.com/firba1/complete"
	"os"
	"os/exec"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("usage:\t%v COMMAND", os.Args[0])
		return
	}

	cmdname := os.Args[1]

	cmdpath, err := exec.LookPath(cmdname)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	linenoise.SetCompletionHandler(
		func(args string) (completions []string) {
			completions = complete.Bash(cmdname + " " + args)
			for i, v := range completions {
				completions[i] = v[len(cmdname)+1:]
			}
			return
		},
	)

	fmt.Println("shell for", cmdpath)

	for {
		line, err := linenoise.Line(fmt.Sprintf(">%v ", cmdname))
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		linenoise.AddHistory(line)

		args := strings.Fields(line)

		cmd := exec.Command(cmdname, args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		err = cmd.Run()
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
	}
}
