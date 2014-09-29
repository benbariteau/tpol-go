package main

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"
)

func ps1(cmdname string) (promptStr string) {
	if cmdname == "git" {
		cmd := exec.Command("bash", "-l", "-c", "__git_ps1")
		gitPs1Stdout, err := cmd.StdoutPipe()
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		if err := cmd.Start(); err != nil {
			fmt.Println(err.Error())
			return
		}

		ps1Bytes, err := ioutil.ReadAll(gitPs1Stdout)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		promptStr = strings.TrimSpace(string(ps1Bytes))

		if err := cmd.Wait(); err != nil {
			fmt.Println(err.Error())
			return
		}
	}
	return
}
