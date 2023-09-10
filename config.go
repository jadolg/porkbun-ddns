package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Record struct {
	Domain      string `yaml:"domain"`
	Host        string `yaml:"host,omitempty"`
	IpV6        bool   `yaml:"ipv6"`
	IpV4        bool   `yaml:"ipv4"`
	Credentials string `yaml:"credentials"`
}

type PorkbunCredentials struct {
	PorkbunAPIKey    string `yaml:"porkbun_api_key"`
	PorkbunSecretKey string `yaml:"porkbun_secret_key"`
}

type MetricsConfig struct {
	Enabled bool `yaml:"enabled"`
	Port    int  `yaml:"port"`
}

type configuration struct {
	Records               []Record                      `yaml:"records"`
	UpdateIntervalMinutes int                           `yaml:"update_interval_minutes"`
	PorkbunCredentials    map[string]PorkbunCredentials `yaml:"credentials"`
	Metrics               MetricsConfig                 `yaml:"metrics"`
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

	for _, credential := range c.PorkbunCredentials {
		if credential.invalid() {
			return configuration{}, fmt.Errorf("invalid credentials configuration detected %+v", credential)
		}
	}

	for _, record := range c.Records {
		if record.invalid() {
			return configuration{}, fmt.Errorf("invalid record configuration detected %+v", record)
		}
		if _, found := c.PorkbunCredentials[record.Credentials]; !found {
			return configuration{}, fmt.Errorf("invalid record credentials detected %+v", record)
		}
	}

	return c, nil
}

func (r Record) invalid() bool {
	return (!r.IpV4 && !r.IpV6) || r.Domain == "" || r.Credentials == ""
}

func (c PorkbunCredentials) invalid() bool {
	return c.PorkbunAPIKey == "" || c.PorkbunSecretKey == ""
}

func ipNeeded(c configuration) (ipv4needed, ipv6needed bool) {
	for _, record := range c.Records {
		if record.IpV4 {
			ipv4needed = true
		}
		if record.IpV6 {
			ipv6needed = true
		}
		if ipv4needed && ipv6needed {
			return
		}
	}
	return
}
