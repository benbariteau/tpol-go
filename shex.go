package main

import (
	"bufio"
	"fmt"
	"github.com/GeertJohan/go.linenoise"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
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

	funcMap := bashCompletionFuncs()

	bashCompletionFunc, ok := funcMap[cmdname]
	var completionFunc linenoise.CompletionHandler = linenoise.DefaultCompletionHandler
	if ok {
		completionFunc = func(args string) []string {
			tmpfile, err := ioutil.TempFile("", cmdname)
			if err != nil {
				fmt.Println("Unable to create temporary file: ", err.Error())
				return []string{}
			}
			tmpfile.Write([]byte(print_completions_src))
			line := fmt.Sprintf("%v %v", cmdname, args)
			completionsCmd := exec.Command(
				"env",
				fmt.Sprintf("COMP_WORDS=(%v)", line),
				"COMP_CWORD=1",
				"bash", "-c",
				fmt.Sprint(
					fmt.Sprintf("source %v;", tmpfile.Name()),
					fmt.Sprintf(
						"source /usr/local/etc/bash_completion; %v; ",
						bashCompletionFunc,
					),
					"__print_completions;",
				),
			)

			stdout, err := completionsCmd.StdoutPipe()
			if err != nil {
				fmt.Println(err.Error())
			}

			completionsCmd.Start()

			bytes, _ := ioutil.ReadAll(stdout)

			fmt.Println(string(bytes))
			return []string{}
		}
	}
	linenoise.SetCompletionHandler(completionFunc)

	fmt.Println("shell for", cmdpath)

	for {
		line, err := linenoise.Line(fmt.Sprintf(">%v ", cmdname))
		if err != nil {
			fmt.Println(err.Error())
			return
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

const print_completions_src = `
__print_completions() {
    for ((i=0;i<${#COMPREPLY[*]};i++))
        do echo $(COMPREPLY[i]}
    done
}
`
