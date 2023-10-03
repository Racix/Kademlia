package d7024e

import (
	"os"
	"os/exec"
	"bytes"
	"testing"
)

func TestCli(t *testing.T) {
	var kad Kademlia
	var stdin bytes.Buffer
	// input := []string{
	// 	"make\nput some data\nput\nget some\nget\nget some some\nexit\n",
	// 	"exit",}

	input := "make\nput some data\nput\nget some\nget\nget some some\nexit\n"

	defer stdin.Reset()

	if os.Getenv("CLI") == "1" {
		stdin.Write([]byte(input))
		kad.Cli(&stdin)	
		return
	}

	// Run the test in a subprocess
	cmd := exec.Command(os.Args[0], "-test.run=TestCli")
	cmd.Env = append(os.Environ(), "CLI=1")
	err := cmd.Run()
	 
	// Cast the error as *exec.ExitError and compare the result
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
        return
    }
    t.Fatalf("process ran with err %v, want exit status 1", err)
}

func TestCli2(t *testing.T) {
	var kad Kademlia
	var stdin bytes.Buffer
	// input := []string{
	// 	"make\nput some data\nput\nget some\nget\nget some some\nexit\n",
	// 	"exit",}

	input := "exit"

	defer stdin.Reset()

	if os.Getenv("CLI") == "1" {
		stdin.Write([]byte(input))
		kad.Cli(&stdin)	
		return
	}

	// Run the test in a subprocess
	cmd := exec.Command(os.Args[0], "-test.run=TestCli2")
	cmd.Env = append(os.Environ(), "CLI=1")
	err := cmd.Run()
	 
	// Cast the error as *exec.ExitError and compare the result
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
        return
    }
    t.Fatalf("process ran with err %v, want exit status 1", err)
}
