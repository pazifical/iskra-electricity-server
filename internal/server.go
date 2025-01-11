package internal

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/pazifical/iskra-electricity-server/iskra"
)

type IskraElectricityServer struct {
	port    int
	mux     *http.ServeMux
	monitor *iskra.ElectricityMonitor
}

func NewIskraElectricityServer(port int, monitor *iskra.ElectricityMonitor) IskraElectricityServer {
	server := IskraElectricityServer{
		port:    port,
		mux:     http.NewServeMux(),
		monitor: monitor,
	}

	server.mux.HandleFunc("/api/electricity/cumulated", server.getCurrentReading)

	return server
}

func (ies *IskraElectricityServer) getCurrentReading(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(ies.monitor.CurrentReading)
	if err != nil {
		log.Printf("ERROR: serving current reading GET request: %v", err)
	}
}

func (ies *IskraElectricityServer) Start() error {
	go ies.monitor.Start()

	log.Printf("INFO: staring server on port %d", ies.port)
	return http.ListenAndServe(fmt.Sprintf(":%d", ies.port), ies.mux)
}
