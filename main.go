package main

import (
	"context"
	"flag"
	"time"

	log "github.com/sirupsen/logrus"
)

func main() {
	var configFile string
	var logLevel string
	var runOnce bool
	flag.StringVar(&configFile, "config", "./config.yaml", "Path to config file")
	flag.StringVar(&logLevel, "loglevel", "info", "Log level")
	flag.BoolVar(&runOnce, "run-once", false, "disables the periodic update")
	flag.Parse()

	err := setLogLevel(logLevel)
	if err != nil {
		log.Fatal(err)
	}

	c, err := getConfig(configFile)
	if err != nil {
		log.Fatal(err)
	}

	updateRecords(c)

	if !runOnce {
		tick := time.Tick(time.Duration(c.UpdateIntervalMinutes) * time.Minute)
		for range tick {
			updateRecords(c)
		}
	}
}

func setLogLevel(logLevel string) error {
	parsedLogLevel, err := log.ParseLevel(logLevel)
	if err != nil {
		return err
	}
	log.SetLevel(parsedLogLevel)
	return err
}

func updateRecords(c configuration) {
	log.Info("Updating records")
	ipv4address, ipv6address := getIpAddresses(c, wtfismyipClient{})

	ctx := context.Background()

	for _, record := range c.Records {
		client, err := getPorkbunClient(c, record.Credentials)
		if err != nil {
			log.Error(err)
			continue
		}
		updateRecord(ctx, record, client, ipv4address, ipv6address)
	}
	log.Info("Records updated")
}
