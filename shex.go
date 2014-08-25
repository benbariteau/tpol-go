package main

import (
	"fmt"
	"os"
	"os/exec"
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

	fmt.Println("shell for ", cmdpath)
}
