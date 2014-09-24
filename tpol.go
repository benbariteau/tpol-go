package main

import (
	"fmt"
	"github.com/GeertJohan/go.linenoise"
	"github.com/firba1/complete"
	"io/ioutil"
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
		ps1 := ""
		if cmdname == "git" {
			cmd := exec.Command("bash", "-l", "-c", "__git_ps1")
			gitPs1Stdout, err := cmd.StdoutPipe()
			if err != nil {
				fmt.Println(err.Error())
			}

			if err := cmd.Start(); err != nil {
				fmt.Println(err.Error())
			}

			ps1Bytes, err := ioutil.ReadAll(gitPs1Stdout)
			if err != nil {
				fmt.Println(err.Error())
			}
			ps1 = strings.TrimSpace(string(ps1Bytes))

			if err := cmd.Wait(); err != nil {
				fmt.Println(err.Error())
			}
		}
		line, err := linenoise.Line(fmt.Sprintf("%v>%v ", ps1, cmdname))
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		linenoise.AddHistory(line)

		// escape to shell
		if len(line) >= 1 && line[0] == '!' {
			cmdFields := strings.Fields(line[1:])
			cmd := exec.Command(cmdFields[0], cmdFields[1:]...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Stdin = os.Stdin
			err = cmd.Run()
			if err != nil {
				fmt.Println(err.Error())
			}
			continue
		}

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
