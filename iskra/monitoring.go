package iskra

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/pazifical/iskra-electricity-server/internal/logging"
	"github.com/pazifical/iskra-electricity-server/internal/types"
	"github.com/pazifical/iskra-electricity-server/sml"
)

// Values over this value are really unlikely and wont be sent to the database
const maxValKWh = 100000

type ElectricityMonitor struct {
	CurrentReading  types.EnergyReading
	readoutInterval int
	currentError    error
}

func NewElectricityMonitor(readoutInterval int) ElectricityMonitor {
	return ElectricityMonitor{
		readoutInterval: readoutInterval,
	}
}

func (em *ElectricityMonitor) Start() {
	logging.Info("staring energy monitoring service")

	for {
		err := configureSerialPort()
		if err != nil {
			logging.Error(err.Error())
		} else {
			break
		}
		time.Sleep(time.Second * 10)
	}

	sleepTime := time.Duration(em.readoutInterval * 1000000000)
	for {
		sensorReadout, ok := readoutSensor()
		if !ok {
			logging.Info(fmt.Sprintf("sleeping for %v minutes.", sleepTime.Minutes()))
			time.Sleep(sleepTime)
			continue
		}

		reading, err := sml.ParseSensorReadout(sensorReadout)
		if err != nil {
			reading = types.EnergyReading{}
			em.currentError = err
			continue
		}

		logging.Info(fmt.Sprintf("energy reading: %v", reading))
		if reading.Value < maxValKWh {
			em.CurrentReading = reading
		} else {
			em.currentError = fmt.Errorf("reading (%f) > max value (%d)", reading.Value, maxValKWh)
			logging.Error(fmt.Sprintf("%f exceeds the max value %d", reading.Value, maxValKWh))
		}

		logging.Info(fmt.Sprintf("sleeping for %v minutes.\n", sleepTime.Minutes()))
		time.Sleep(sleepTime)
	}
}

// Configures the serial port so that a readout is possible.
// The configuration is necessary after a system reboot
func configureSerialPort() error {
	logging.Info("trying to configure serial port")

	out, err := exec.Command(
		"stty",
		"-F",
		"/dev/ttyUSB0",
		"1:0:8bd:0:3:1c:7f:15:4:5:1:0:11:13:1a:0:12:f:17:16:0:0:0:0:0:0:0:0:0:0:0:0:0:0:0:0",
	).Output()

	logging.Debug(string(out))

	if err != nil {
		return fmt.Errorf("cannot configure serial port: %v", err)
	}
	return nil
}
