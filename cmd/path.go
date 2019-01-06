package cmd

import (
	"errors"
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v2"
)

//ServerConfig contains map of services directories
type ServerConfig struct {
	Dir      string
	Services map[string]map[string]string `yaml:"services"`
	Port     int                          `yaml:"port"`
}

//Init server config
func (c *ServerConfig) Init() error {
	if err := c.SetDir(); err != nil {
		return err
	}
	if err := c.ReadConfig(); err != nil {
		return err
	}
	if err := c.CheckDirectories(); err != nil {
		return err
	}
	return nil
}

//SetDir reads current pwd
func (c *ServerConfig) SetDir() error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	c.Dir = dir
	return nil
}

//ReadConfig reads services from yml file
func (c *ServerConfig) ReadConfig() error {
	data, err := ioutil.ReadFile(c.Dir + "/flaka-ci.yml")
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal([]byte(data), &c); err != nil {
		return err
	}
	return nil
}

//CheckDirectories check if directories exist
func (c *ServerConfig) CheckDirectories() error {
	for _, options := range c.Services {
		if options["path"] == "" {
			return errors.New("Error in flaka-ci.yml: path option is required")
		}
		if _, err := os.Stat(c.Dir + "/" + options["path"]); os.IsNotExist(err) {
			return errors.New("flaka-ci.yml: Could not find directory " + options["path"] + " in " + c.Dir)
		}
	}
	return nil
}
