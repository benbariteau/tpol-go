package main

import (
    "fmt"
    "os"
    "os/exec"
)

func main() {
    cmd := exec.Command("bash", "-c", "source /etc/bash_completion; complete -p")
    
}
