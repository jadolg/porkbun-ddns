package main

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/nrdcg/porkbun"
	log "github.com/sirupsen/logrus"
)

func createOrUpdateRecord(ctx context.Context, client *porkbun.Client, recordID string, domain, host, recordType, address string) error {
	if recordID != "" {
		recordIDint, err := strconv.Atoi(recordID)
		if err != nil {
			return err
		}
		log.WithFields(log.Fields{
			"record id":   recordID,
			"domain":      domain,
			"host":        host,
			"record type": recordType,
			"address":     address,
		}).Info("Updating record")
		return client.EditRecord(ctx, domain, recordIDint, porkbun.Record{
			Name:    host,
			Type:    recordType,
			Content: address,
			Notes:   "DDNS",
		})
	} else {
		log.WithFields(log.Fields{
			"record id":   recordID,
			"domain":      domain,
			"host":        host,
			"record type": recordType,
			"address":     address,
		}).Info("Creating record")
		_, err := client.CreateRecord(ctx, domain, porkbun.Record{
			Name:    host,
			Type:    recordType,
			Content: address,
			Notes:   "DDNS",
		})
		return err
	}
}

func getRecords(ctx context.Context, domain, host string, client *porkbun.Client) (*porkbun.Record, *porkbun.Record, error) {
	records, err := client.RetrieveRecords(ctx, domain)
	if err != nil {
		return nil, nil, err
	}
	var ipv4Record porkbun.Record
	var ipv6Record porkbun.Record

	for _, record := range records {
		if record.Name == fmt.Sprintf("%s.%s", host, domain) || (host == "" && record.Name == domain) {
			if record.Type == "A" {
				ipv4Record = record
			} else if record.Type == "AAAA" {
				ipv6Record = record
			}

		}
	}
	return &ipv4Record, &ipv6Record, nil
}

func getPorkbunClients(credentials map[string]PorkbunCredentials) (map[string]*porkbun.Client, error) {
	clients := make(map[string]*porkbun.Client)

	for key, credential := range credentials {
		client, err := getPorkbunClient(credential, key)
		if err != nil {
			connectionErrorsTotal.Inc()
			log.Errorf("Error getting client for credentials '%s': %v", key, err)
			continue
		}
		clients[key] = client
	}
	return clients, nil
}

func getPorkbunClient(credentials PorkbunCredentials, credentialsName string) (*porkbun.Client, error) {
	client := porkbun.New(credentials.PorkbunSecretKey, credentials.PorkbunAPIKey)

	ctx := context.Background()

	yourIP, err := client.Ping(ctx)
	if err != nil {
		return nil, err
	}

	log.Debugf("Connected to porkbun from %s using credentials %s", yourIP, credentialsName)
	return client, nil
}

func updateRecord(ctx context.Context, record Record, client *porkbun.Client, ipv4address string, ipv6address string) (resultError error) {
	ipv4Record, ipv6Record, err := getRecords(ctx, record.Domain, record.Host, client)
	if err != nil {
		resultError = errors.Join(resultError, err)
	}
	if record.IpV4 && ipv4Record != nil && ipv4address != "" && ipv4Record.Content != ipv4address {
		err = createOrUpdateRecord(ctx, client, ipv4Record.ID, record.Domain, record.Host, "A", ipv4address)
		if err != nil {
			resultError = errors.Join(resultError, err)
		} else {
			updateSuccessTotal.WithLabelValues(record.Host, record.Domain, "A").Inc()
		}
	} else {
		log.WithFields(log.Fields{
			"host":   record.Host,
			"domain": record.Domain,
		}).Debug("Ipv4 update not required")
	}
	if record.IpV6 && ipv6Record != nil && ipv6address != "" && ipv6Record.Content != ipv6address {
		err = createOrUpdateRecord(ctx, client, ipv6Record.ID, record.Domain, record.Host, "AAAA", ipv6address)
		if err != nil {
			resultError = errors.Join(resultError, err)
		} else {
			updateSuccessTotal.WithLabelValues(record.Host, record.Domain, "AAAA").Inc()
		}
	} else {
		log.WithFields(log.Fields{
			"host":   record.Host,
			"domain": record.Domain,
		}).Debug("Ipv6 update not required")
	}
	return
}
