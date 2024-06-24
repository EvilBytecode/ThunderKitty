package FactoryReset

import "os/exec"

func Disable() {
    exec.Command("reagentc.exe", "/disable").Run()
}