package config

import (
	"gopkg.in/gcfg.v1"
	"log"
)

type MainConfig struct {
	App      AppConfig
	Maps     MapsConfig
	Database DBConfig
}

type AppConfig struct {
	Limit int
}

type MapsConfig struct {
	Key string
}

type DBConfig struct {
	DSN string
}

var mainConfig *MainConfig

func ReadConfig(path, fileName string) bool {
	mainConfig = &MainConfig{}
	err := gcfg.ReadFileInto(mainConfig, path+fileName+".ini")

	if nil != err {
		log.Println("error reading config", fileName, err)
		return false
	}

	return true
}

func Get() *MainConfig {
	return mainConfig
}
