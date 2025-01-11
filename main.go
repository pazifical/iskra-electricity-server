package main

import (
	"log"
	"os"
	"strconv"

	"github.com/pazifical/iskra-electricity-server/internal"
	"github.com/pazifical/iskra-electricity-server/iskra"
)

var readoutInterval = 60
var port = 8080

func init() {
	var err error

	envPort := os.Getenv("IES_PORT")
	if envPort != "" {
		port, err = strconv.Atoi(envPort)
		if err != nil {
			log.Fatal(err)
		}
	}

	envReadoutInterval := os.Getenv("IES_READOUT_INTERVAL")
	if envReadoutInterval != "" {
		readoutInterval, err = strconv.Atoi(envReadoutInterval)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	monitor := iskra.NewElectricityMonitor(readoutInterval)

	server := internal.NewIskraElectricityServer(port, &monitor)

	err := server.Start()
	if err != nil {

	}
}
