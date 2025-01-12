package iskra

import (
	"os/exec"
	"strings"

	"github.com/pazifical/iskra-electricity-server/internal/logging"
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
	var err error
	binToHexCmd.Stdin, err = catSerialCmd.StdoutPipe()
	if err != nil {
		logging.Error(err.Error())
		return "", false
	}

	binToHexCmd.Stdout = stdoutBuilder

	logging.Info("starting sensor readout")
	err = binToHexCmd.Start()
	if err != nil {
		logging.Error(err.Error())
		return "", false
	}

	err = catSerialCmd.Run()
	if err != nil {
		logging.Error(err.Error())
		return "", false
	}
	err = binToHexCmd.Wait()
	if err != nil {
		logging.Error(err.Error())
		return "", false
	}

	logging.Info("finished sensor readout")

	readout := stdoutBuilder.String()
	if len(readout) == 0 {
		logging.Error("readout is empty.")
		return "", false
	}
	return readout, true
}
