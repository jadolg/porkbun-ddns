package main

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

const (
	// DefaultWTFIsMyIPV4Endpoint is the default endpoint for IPv4 address retrieval
	DefaultWTFIsMyIPV4Endpoint = "https://ipv4.wtfismyip.com/json"
	// DefaultWTFIsMyIPV6Endpoint is the default endpoint for IPv6 address retrieval
	DefaultWTFIsMyIPV6Endpoint = "https://ipv6.wtfismyip.com/json"
)

type Record struct {
	Domain      string `yaml:"domain"`
	Host        string `yaml:"host,omitempty"`
	IpV6        bool   `yaml:"ipv6"`
	IpV4        bool   `yaml:"ipv4"`
	Notes       string `yaml:"notes,omitempty"`
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

type WTFIsMyIPConfig struct {
	V4Endpoint string `yaml:"v4"`
	V6Endpoint string `yaml:"v6"`
}

type configuration struct {
	Records               []Record                      `yaml:"records"`
	UpdateIntervalMinutes int                           `yaml:"update_interval_minutes"`
	PorkbunCredentials    map[string]PorkbunCredentials `yaml:"credentials"`
	Metrics               MetricsConfig                 `yaml:"metrics"`
	Timeout               time.Duration                 `yaml:"timeout"`
	WTFIsMyIP             WTFIsMyIPConfig               `yaml:"wtfismyip"`
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

	if c.Timeout == 0 {
		c.Timeout = 10 * time.Second
	}

	if c.UpdateIntervalMinutes == 0 {
		c.UpdateIntervalMinutes = 5
	}

	if c.WTFIsMyIP.V4Endpoint == "" {
		c.WTFIsMyIP.V4Endpoint = DefaultWTFIsMyIPV4Endpoint
	}

	if c.WTFIsMyIP.V6Endpoint == "" {
		c.WTFIsMyIP.V6Endpoint = DefaultWTFIsMyIPV6Endpoint
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
