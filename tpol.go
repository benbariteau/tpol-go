package main

import (
	"fmt"
	"github.com/GeertJohan/go.linenoise"
	"github.com/firba1/complete"
	"os"
	"os/exec"
	"os/signal"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("usage:\t%v COMMAND\n", os.Args[0])
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
			if len(args) > 1 && args[0] == '!' {
				cmd := args[1:]
				completions = complete.Complete(cmd)
				for i, v := range completions {
					completions[i] = "!" + v
				}
				return completions
			}
			completions = complete.Complete(cmdname + " " + args)
			for i, v := range completions {
				completions[i] = v[len(cmdname)+1:]
			}
			return
		},
	)

	ps := NewPromptStringer(PromptStringMapping{"git", "__git_ps1"})

	fmt.Println("shell for", cmdpath)

	for {
		promptStr := ps.PromptString(cmdname)
		line, err := linenoise.Line(fmt.Sprintf("%v>%v ", promptStr, cmdname))
		if err != nil {
			fmt.Println(err.Error())
			break
		}
		linenoise.AddHistory(line) // add history

		var cmd *exec.Cmd
		if len(line) >= 1 && line[0] == '!' { // escape to shell
			cmdFields := strings.Fields(line[1:])
			cmd = exec.Command(cmdFields[0], cmdFields[1:]...)
		} else if strings.TrimSpace(line) == "exit" { // exit
			break
		} else { // regular subcommand
			args := strings.Fields(line)
			cmd = exec.Command(cmdname, args...)
		}

		// hook up expecting tty stuff
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin

		// handle interrupts
		cancel := catchAndPassSignal(cmd, os.Interrupt)

		err = cmd.Run()
		close(cancel)

		if err != nil {
			fmt.Println(err.Error())
			continue
		}
	}
}

/*
catchAndPassSignal catches the given signals and passes them to the process of the given command

catchAndPassSignal can be canceled by closing the cancel channel that it returns
*/
func catchAndPassSignal(cmd *exec.Cmd, signals ...os.Signal) (cancel chan int) {
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, signals...)

	cancel = make(chan int)

	go func() {
		select {
		case <-cancel:
			return
		case sig := <-sigint:
			cmd.Process.Signal(sig)
		}
	}()
	return
}
