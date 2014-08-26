package main

import (
	"bufio"
	"fmt"
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

	fmt.Println("shell for", cmdpath)

	stdin := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf(">%v ", cmdname)
		line, err := stdin.ReadString('\n')
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		line = line[:len(line)-1] // remove trailing newline
		args := strings.Split(line, " ")
		cmd := exec.Command(cmdname, args...)
		out, err := cmd.Output()
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		fmt.Print(string(out))
	}
}
