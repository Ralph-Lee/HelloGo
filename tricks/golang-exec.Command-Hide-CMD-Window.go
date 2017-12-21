package main

func main() {
	cmd := exec.Command("explorer", s.GetLocalUrl())
	if err := cmd.Run(); err != nil {
		logger.Error("explorer start failed!", err.Error())
	}
}

// C:\> go build -ldflags "-s -w -H=windowsgui"
