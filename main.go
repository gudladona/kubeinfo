package main

import (
	"flag"
	"fmt"
	"github.com/gudladona87/kubeinfo/config"
	"log"
	"net/http"
)

func main() {
	log.Print("Starting KubeInfo Service")

	//Read configuration from file
	configfile := flag.String("config", "/etc/kubeinfo/config.cfg", "/path/to/configfile")
	flag.Parse()

	err := config.ReadConfigFromFile(*configfile)
	if err != nil {
		log.Fatalf("error reading from config file: %v", err)
	}

	// Start HTTP Server
	err = Start()
	if err != nil {
		log.Fatalf("error starting http server: %v", err)
	}
}

//Start will start up the web server for serving HTTP requests
func Start() error {
	router, err := registerHandlers()
	if err != nil {
		return err
	}
	return http.ListenAndServe(fmt.Sprintf(":%d", config.ListenPort()), router)
}
