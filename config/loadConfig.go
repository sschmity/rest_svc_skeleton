package config

import (
	"os"
	"log"
	"encoding/json"
	"io/ioutil"
)

const ConfigFileName = "config.json"

// ConfigFile struct from config.json file
type ConfigFileStruct struct {
	HttpServerConfig	HttpServerConfigStruct `json:"http.server.config"`
	LogConfig			LogConfigStruct `json:"log.config"`
}

type HttpServerConfigStruct struct {
	Name	string `json:"name"`
	Port	string `json:"port"`
}

type LogConfigStruct struct {
	FilePath	string `json:"file.path"`
	FileName	string `json:"file.name"`
}

var ConfigFileInstance *ConfigFileStruct

func LoadConfigFile(aConfigFileName string) {
	jsonFile, err := os.Open(aConfigFileName)
	defer jsonFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	// read our opened jsonFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &ConfigFileInstance)

	log.Println(ConfigFileInstance)
}
