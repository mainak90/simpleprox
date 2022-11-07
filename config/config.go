package config

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
)

type ProxyConfig struct {
	ListenHost string
	BackendHost string
}

func New() *ProxyConfig {
	fullpath, ok := os.LookupEnv("CONFIG_PATH")
	if !ok {
		fullpath = "config.json"
	}
	log.WithFields(log.Fields{"Controller": "config"}).Info("Loading config from file: ", fullpath)
	res := ProxyConfig{}
	file, err := os.Open(fullpath)
	if err != nil {
		log.WithFields(log.Fields{"Controller": "config", "Error": "ENOENT"}).Error("Error opening file: ", err)
		return nil
	}
	defer file.Close()
	buffer, err := ioutil.ReadAll(file)
	if err != nil {
		log.WithFields(log.Fields{"Controller": "config", "Error": "Empty"}).Error("File empty: ", err)
		return nil
	}
	err = json.Unmarshal(buffer, &res)
	if err != nil {
		log.WithFields(log.Fields{"Controller": "config", "Error": "Unmarshallable"}).Error("Cannot marshal json into struct: ", err)
		return nil
	}
	log.WithFields(log.Fields{"Controller": "config"}).Info("Startup config loaded.")
	return &res
}