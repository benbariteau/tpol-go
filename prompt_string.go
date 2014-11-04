package main

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"
)

type PromptStringer map[string]string

type PromptStringMapping struct {
	Command string
	Func    string
}

func NewPromptStringer(mappings ...PromptStringMapping) PromptStringer {
	ps := PromptStringer(make(map[string]string))
	for _, mapping := range mappings {
		ps[mapping.Command] = mapping.Func
	}
	return ps
}

func (ps PromptStringer) PromptString(cmdname string) (promptStr string) {
	if psFunc := ps[cmdname]; psFunc != "" {
		cmd := exec.Command("bash", "-l", "-c", psFunc)
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
			return
		}
	}
	return
}
