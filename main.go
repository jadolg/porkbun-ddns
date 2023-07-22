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
	flag.StringVar(&configFile, "config", "./config.yaml", "Path to config file")
	flag.StringVar(&logLevel, "loglevel", "info", "Log level")
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

	tick := time.Tick(time.Duration(c.UpdateIntervalMinutes) * time.Minute)
	for range tick {
		updateRecords(c)
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
	ipv4address, err := getIpAddress(false)
	if err != nil {
		log.Error(err)
	}
	log.Debugf("Current host ipv4: %s", ipv4address)
	ipv6address, err := getIpAddress(true)
	if err != nil {
		log.Error(err)
	}
	log.Debugf("Current host ipv6: %s", ipv6address)

	ctx := context.Background()

	for _, record := range c.Records {
		client, err := getPorkbunClient(c, record.Credentials)
		if err != nil {
			log.Error(err)
			continue
		}
		ipv4Record, ipv6Record, err := getRecords(ctx, record.Domain, record.Host, client)
		if err != nil {
			log.Error(err)
		}
		if record.IpV4 && ipv4Record != nil && ipv4address != "" && ipv4Record.Content != ipv4address {
			err = createOrUpdateRecord(ctx, client, ipv4Record.ID, record.Domain, record.Host, "A", ipv4address)
		} else {
			log.Debugf("Ipv4 update not required for %s.%s", record.Host, record.Domain)
		}
		if record.IpV6 && ipv6Record != nil && ipv6address != "" && ipv6Record.Content != ipv6address {
			err = createOrUpdateRecord(ctx, client, ipv6Record.ID, record.Domain, record.Host, "AAAA", ipv6address)
		} else {
			log.Debugf("Ipv6 update not required for %s.%s", record.Host, record.Domain)
		}
	}

	if err != nil {
		log.Error(err)
	}
}
