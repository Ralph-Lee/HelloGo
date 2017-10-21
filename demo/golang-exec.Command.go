/**
********************************************************************************* 
* http://www.darrencoxall.com/golang/executing-commands-in-go/
* https://nathanleclaire.com/blog/2014/12/29/shelled-out-commands-in-golang/
********************************************************************************* 
* There are three different types of command execution within applications.
*
* [1][Plain output] - You expect the command to execute but all you want from it is the output.
*
* [2][Exit codes] - You want to discard any output but assert the exit code.
*
* [3][Long running processes] - You want to spawn sub-processes.
********************************************************************************* 
**/

package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"strings"
	"syscall"
	"io/ioutil"
	"bytes"
	"time"
)

var oscmdfile = []byte(`
C:\Windows\System32\cmd.exe
C:\Windows\System32\PowerShell\1.0\powershell.exe
`)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func genCommandFile() {
	my, err := user.Current()
	check(err)
	err = ioutil.WriteFile(my.HomeDir + string(os.PathSeparator) + "donotrun.ps1", oscmdfile, 0755)
	check(err)
}

func printOutput(outs []byte) {
	if len(outs) > 0 {
		fmt.Printf("exec.Output: %s\n", string(outs))
	}
}

func printCommand(cmd *exec.Cmd) {
	fmt.Printf("exec.Executing: %s\n", strings.Join(cmd.Args, " "))
}

func printError(err error) {
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("exec.Error: %s\n", err.Error()))
	}
}

func demoCollectingOutputEasy() {
	//
	// Synchronous
	// Combine stdout and stderr
	//
	cmd := exec.Command("echo", "Called from Go!")
	printCommand(cmd)
	output, err := cmd.CombinedOutput()
	printError(err)
	printOutput(output)
}

func demoCollectingOutputDifBuffers() {
	//
	// Synchronous
	//
	cmd := exec.Command("go", "version")
	cmdOutput := &bytes.Buffer{} // Stdout buffer
	cmd.Stdout = cmdOutput // Attach buffer to command
	printCommand(cmd) // Execute command
	err := cmd.Run() // will wait for command to return
	printError(err)
	printOutput(cmdOutput.Bytes()) // Only output the commands stdout
}

func demoRetrievingExitCode() {
	//
	// Synchronous
	//
	cmd := exec.Command("ls", "/imaginary/directory")
	var waitStatus syscall.WaitStatus
	if err := cmd.Run(); err != nil {
		printError(err)
		// Did the command fail because of an unsuccessful exit code
		if exitError, ok := err.(*exec.ExitError); ok {
			waitStatus = exitError.Sys().(syscall.WaitStatus)
			printOutput([]byte(fmt.Sprintf("%d", waitStatus.ExitStatus())))
		}
	} else {
		// Command was successful
		waitStatus = cmd.ProcessState.Sys().(syscall.WaitStatus)
		printOutput([]byte(fmt.Sprintf("%d", waitStatus.ExitStatus())))
	}
}

func demoLongRunningProcesses(){
	//
	// Asynchronously
	//
	cmd := exec.Command("cat", "/dev/random")
	randomBytes := &bytes.Buffer{}
	cmd.Stdout = randomBytes

	// Start command asynchronously
	err := cmd.Start()
	printError(err)

	// Create a ticker that outputs elapsed time
	ticker := time.NewTicker(time.Second)
	go func(ticker *time.Ticker) {
		now := time.Now()
		for _ = range ticker.C {
			printOutput(
				[]byte(fmt.Sprintf("%s", time.Since(now))),
			)
		}
	}(ticker)

	// Create a timer that will kill the process
	timer := time.NewTimer(time.Second * 4)
	go func(timer *time.Timer, ticker *time.Ticker, cmd *exec.Cmd) {
		for _ = range timer.C {
			err := cmd.Process.Signal(os.Kill)
			printError(err)
			ticker.Stop()
		}
	}(timer, ticker, cmd)

	// Only proceed once the process has finished
	cmd.Wait()
	printOutput(
		[]byte(fmt.Sprintf("%d bytes generated!", len(randomBytes.Bytes()))),
	)
}

func main() {
	//genCommandFile ()
	cmd := exec.Command("C://Windows//System32//WindowsPowerShell//v1.0//powershell.exe", "Start-Process PowerShell -Verb runas")
}
