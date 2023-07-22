package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Record struct {
	Domain string `yaml:"domain"`
	Host   string `yaml:"host"`
	IpV6   bool   `yaml:"ipv6"`
	IpV4   bool   `yaml:"ipv4"`
}

type configuration struct {
	PorkbunAPIKey         string
	PorkbunSecretKey      string
	Records               []Record `yaml:"records"`
	UpdateIntervalMinutes int      `yaml:"update_interval_minutes"`
}

func getConfig(configFile string) (configuration, error) {
	yamlFile, err := os.ReadFile(configFile)
	if err != nil {
		return configuration{}, err
	}
	c := configuration{}
	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		return configuration{}, err
	}

	if c.PorkbunAPIKey == "" || c.PorkbunSecretKey == "" {
		return configuration{}, fmt.Errorf("the Porkbun API key and Secret key should be provided in the configuration")
	}

	for _, record := range c.Records {
		if record.invalid() {
			return configuration{}, fmt.Errorf("invalid record configuration detected %+v", record)
		}
	}

	return c, nil
}

func (r Record) invalid() bool {
	return (!r.IpV4 && !r.IpV6) || r.Domain == "" || r.Host == ""
}
