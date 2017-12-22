// go build -ldflags "-s -w -H=windowsgui"
package main

import "os/exec"

func main() {
	cmd := exec.Command("explorer", s.GetLocalUrl())
	if err := cmd.Run(); err != nil {
		logger.Error("explorer start failed!", err.Error())
	}
}
