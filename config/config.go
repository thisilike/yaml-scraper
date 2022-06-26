package config

import (
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

var Logger *logrus.Logger
var Err error

type ConfigStruct struct {
	LogLevel string `yaml:"logLevel"`
}

func readConfig() (ConfigStruct, error) {
	data, err := ioutil.ReadFile("config.yml")
	if err != nil {
		fmt.Println("failed accesing config file config.yml")
		return ConfigStruct{}, errors.New("failed accesing config file config.yml")
	}

	config := ConfigStruct{}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		fmt.Println("failed parsing config file config.yml")
		return ConfigStruct{}, errors.New("failed parsing config file config.yml")
	}
	return config, nil
}

func validateConfig(config ConfigStruct) error {
	// validating log levels
	validLogLevels := []string{"TRACE", "DEBUG", "INFO", "WARN", "ERROR", "FATAL", "PANIC"}
	valid := false
	for _, lvl := range validLogLevels {
		if config.LogLevel == lvl {
			valid = true
			break
		}
	}
	if !valid {
		msg := fmt.Sprintf("'%s' is not a valid logLevel. valid log levels are: '%v'",
			config.LogLevel,
			validLogLevels,
		)
		return errors.New(msg)
	}
	return nil
}

func loadConfig(config ConfigStruct) error {
	// initiating logger
	switch config.LogLevel {
	case "TRACE":
		Logger.SetLevel(logrus.TraceLevel)
	case "DEBUG":
		Logger.SetLevel(logrus.DebugLevel)
	case "INFO":
		Logger.SetLevel(logrus.InfoLevel)
	case "WARN":
		Logger.SetLevel(logrus.WarnLevel)
	case "ERROR":
		Logger.SetLevel(logrus.ErrorLevel)
	case "FATAL":
		Logger.SetLevel(logrus.FatalLevel)
	case "PANIC":
		Logger.SetLevel(logrus.PanicLevel)
	default:
		fmt.Printf("invalid logger level: '%s'", config.LogLevel)
		panic("invalid logger level")
	}
	return nil
}

func init() {
	Logger = logrus.New()
	config, err := readConfig()
	if err != nil {
		Err = err
		return
	}

	err = validateConfig(config)
	if err != nil {
		Err = err
		return
	}

	err = loadConfig(config)
	if err != nil {
		Err = err
		return
	}
}
