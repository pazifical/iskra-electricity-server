package iskra

import (
	"log"
	"os/exec"
	"strings"
)

// readoutSensor reads from /dev/ttyUSB0 for 10 seconds
// Since it is serial and the valuable information (energy consumption) is
// right at the beginning, we can timeout the whole command.
// od transforms the binary output to an easy to interpret hex format
func readoutSensor() (string, bool) {
	catSerialCmd := exec.Command("timeout", "10", "cat", "/dev/ttyUSB0")
	binToHexCmd := exec.Command("od", "-tx1")

	// Configuring the piping and how the two commands interact with each other
	stdoutBuilder := new(strings.Builder)
	binToHexCmd.Stdin, _ = catSerialCmd.StdoutPipe()
	binToHexCmd.Stdout = stdoutBuilder

	// Reading from serial for a set time
	log.Println("INFO: Starting sensor readout.")
	_ = binToHexCmd.Start()
	_ = catSerialCmd.Run()
	_ = binToHexCmd.Wait()
	log.Println("INFO: Finished sensor readout.")

	readout := stdoutBuilder.String()
	if len(readout) == 0 {
		log.Printf("ERROR: Readout is empty.")
		return "", false
	}
	return readout, true
}
