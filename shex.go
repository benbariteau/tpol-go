package main

import (
	"bufio"
	"fmt"
	"github.com/GeertJohan/go.linenoise"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"runtime"
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
		completionFunc = func(args string) (completions []string) {
			tmpfile, err := ioutil.TempFile("", cmdname)
			if err != nil {
				fmt.Println("Unable to create temporary file: ", err.Error())
				return
			}
			tmpfile.Write([]byte(print_completions_src))
			line := fmt.Sprintf("%v %v", cmdname, args)
			completionsCmd := exec.Command(
				"env",
				fmt.Sprintf("COMP_CWORD=%v", len(strings.Fields(line))-1),
				fmt.Sprintf("COMP_LINE=%v", line),
				fmt.Sprintf("COMP_POINT=%v", len(line)+1),
				"bash", "-c",
				fmt.Sprint(
					fmt.Sprintf("source %v;", tmpfile.Name()),
					fmt.Sprintf("COMP_WORDS=(%v);", line),
					fmt.Sprintf(
						"source %v; %v; ",
						bashCompletionPath(),
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

			out := bufio.NewReader(stdout)

			argarr := strings.Fields(args)
			prefix := strings.Join(argarr[:len(argarr)-1], " ")
			if prefix != "" {
				prefix = prefix + " "
			}
			for line, err := out.ReadString('\n'); err == nil; line, err = out.ReadString('\n') {
				completions = append(
					completions,
					fmt.Sprintf(
						"%v%v",
						prefix,
						line[:len(line)-1],
					),
				)
			}
			return
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

func bashCompletionFuncs() map[string]string {
	cmd := exec.Command("bash", "-c", fmt.Sprintf("source %v; complete -p", bashCompletionPath()))

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

const bashCompletionPathUnix = "/etc/bash_completion"

func bashCompletionPath() string {
	switch runtime.GOOS {
	case "darwin":
		return brewPrefix() + bashCompletionPathUnix
	default:
		return bashCompletionPathUnix
	}
}

func brewPrefix() string {
	cmd := exec.Command("brew", "--prefix")

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		panic(err.Error())
	}

	cmd.Start()

	stdout := bufio.NewReader(stdoutPipe)
	line, err := stdout.ReadString('\n')
	if err != nil {
		panic(err.Error())
	}

	if err := cmd.Wait(); err != nil {
		return ""
	}

	return line[:len(line)-1]
}

const print_completions_src = `
__print_completions() {
    for ((i=0;i<${#COMPREPLY[*]};i++))
        do echo ${COMPREPLY[i]}
    done
}
`
