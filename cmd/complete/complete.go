package main

import (
    "bufio"
    "fmt"
    "os/exec"
    "regexp"
)

func main() {
    cmd := exec.Command("bash", "-c", "source /etc/bash_completion; complete -p")
    
    stdout, err := cmd.StdoutPipe()
    if err != nil {
        fmt.Println(err.Error())
    }
    scanner := bufio.NewScanner(stdout)

    cmd.Start()

    completionFuncs := make(map[string]string)

    rx := regexp.MustCompile("-F ([^ ]*) (.*)")
    for scanner.Scan() {
        match := rx.FindStringSubmatch(scanner.Text())
        if match != nil {
            funcName := match[1]
            cmdName := match[2]
            completionFuncs[cmdName] = funcName
        }
    }

    for cmdName, funcName := range completionFuncs {
        fmt.Println(cmdName, funcName)
    }

    if err := scanner.Err(); err != nil {
        fmt.Println(err.Error())
    }
}
