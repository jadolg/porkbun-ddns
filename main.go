package main

import (
	"context"
	"flag"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/term"
)

var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
	BuiltBy = "dirty hands"
)

func main() {
	var configFile string
	var logLevel string
	var runOnce bool
	var version bool

	flag.StringVar(&configFile, "config", "./config.yaml", "Path to config file")
	flag.StringVar(&logLevel, "loglevel", "info", "Log level")
	flag.BoolVar(&runOnce, "run-once", false, "disables the periodic update")
	flag.BoolVar(&version, "version", false, "prints version and exits")
	flag.Parse()

	err := configureLogs(logLevel)
	if err != nil {
		log.Fatal(err)
	}

	log.Infof("Version: %s, Commit: %s, Date: %s, BuiltBy: %s", Version, Commit, Date, BuiltBy)
	if version {
		return
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

func configureLogs(logLevel string) error {
	parsedLogLevel, err := log.ParseLevel(logLevel)
	if err != nil {
		return err
	}
	log.SetLevel(parsedLogLevel)
	if !term.IsTerminal(int(os.Stdout.Fd())) {
		log.SetFormatter(&log.JSONFormatter{})
	}
	return err
}

func updateRecords(c configuration) {
	clients, err := getPorkbunClients(c.PorkbunCredentials)
	if err != nil {
		log.Fatal(err)
	}

	log.Info("Updating records")
	ipv4address, ipv6address := getIpAddresses(c, wtfismyipClient{})

	ctx := context.Background()

	for _, record := range c.Records {
		client, exists := clients[record.Credentials]
		if !exists {
			log.Errorf("Credentials not found: %s", record.Credentials)
			continue
		}
		err := updateRecord(ctx, record, client, ipv4address, ipv6address)
		if err != nil {
			log.WithFields(log.Fields{
				"host":        record.Host,
				"domain":      record.Domain,
				"IPv4":        record.IpV4,
				"IPv6":        record.IpV6,
				"credentials": record.Credentials,
				"IPv4Address": ipv4address,
				"IPv6Address": ipv6address,
			}).Error(err)
		}
	}
	log.Info("Records updated")
}
