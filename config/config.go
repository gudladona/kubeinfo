package config

import (
	"sync"
	"io/ioutil"
	"github.com/BurntSushi/toml"
)

var (
	config     *Configuration
	configLock = new(sync.RWMutex)
)

// Configuration encapsulates the user configurable options for the service
type Configuration struct {
	ServerCfg ServerCfg `toml:"server"`
}

// ServerCfg defines the user configurable options for the server itself
type ServerCfg struct {
	ListenPort int    `toml:"listenport"`
	ServerName string `toml:"serverName"`
}

// ReadConfigFromFile is the primary loading function to pull data from the toml config file and sync it.
func ReadConfigFromFile(cfgFile string) error {
	configFile, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		return err
	}
	tempConf := &Configuration{}
	err = toml.Unmarshal(configFile, tempConf)
	if err != nil {
		return err
	}
	configLock.Lock()
	defer configLock.Unlock()
	config = tempConf

	return nil
}

// ListenPort function returns the Port on which the service will be listening.
func ListenPort() int {
	configLock.RLock()
	defer configLock.RUnlock()
	return config.ServerCfg.ListenPort
}
