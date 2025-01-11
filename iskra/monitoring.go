package iskra

import (
	"fmt"
	"log"
	"os/exec"
	"time"

	"github.com/pazifical/iskra-electricity-server/internal/types"
)

// Values over this value are really unlikely and wont be sent to the database
const maxValKWh = 100000

type ElectricityMonitor struct {
	CurrentReading  types.EnergyReading
	readoutInterval int
}

func NewElectricityMonitor(readoutInterval int) ElectricityMonitor {
	return ElectricityMonitor{
		readoutInterval: readoutInterval,
	}
}

func (em *ElectricityMonitor) Start() {
	fmt.Println("ElectricityTotal meter readout and processing.")

	configureSerialPort()

	sleepTime := time.Duration(em.readoutInterval * 1000000000)
	for {
		// Reading out the sensor
		sensorReadout, ok := readoutSensor()
		if !ok {
			log.Printf("INFO: Sleeping for %v minutes.", sleepTime.Minutes())
			time.Sleep(sleepTime)
			continue
		}

		// Processing the readout to get the energy consumption
		// preprocessedReadout := preProcessReadout(sensorReadout)
		// reading, ok := extractConsumptionFromReadout(preprocessedReadout)
		reading, err := ParseSML(sensorReadout)

		if err != nil {
			reading = types.EnergyReading{Error: err}
			continue
		}

		log.Println("Energy reading:", reading)
		if reading.Value < maxValKWh {
			em.CurrentReading = reading
		} else {
			reading.Error = fmt.Errorf("reading (%f) > max value (%d)", reading.Value, maxValKWh)
			em.CurrentReading = reading
			log.Printf("ERROR: %f exceeds the max value %d", reading.Value, maxValKWh)
		}

		// Sleeping for the given time period
		log.Printf("INFO: Sleeping for %v minutes.\n", sleepTime.Minutes())
		time.Sleep(sleepTime)
	}
}

// Configures the serial port so that a readout is possible.
// The configuration is necessary after a system reboot
func configureSerialPort() {
	log.Println("INFO: Trying to configure serial port...")

	out, err := exec.Command(
		"stty",
		"-F",
		"/dev/ttyUSB0",
		"1:0:8bd:0:3:1c:7f:15:4:5:1:0:11:13:1a:0:12:f:17:16:0:0:0:0:0:0:0:0:0:0:0:0:0:0:0:0",
	).Output()

	if err != nil {
		log.Printf("ERROR: Cannot configure serial port: %s", err)
	}
	log.Println(out)
}
