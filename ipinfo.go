package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// IPInfoData is a data representation with the same structure returned by https://ip.r4bbit.net/json
type IPInfoData struct {
	IPAddress string `json:"IPAddress"`
}

type IPInfoClient interface {
	getIpAddress(ipv6 bool) (string, error)
}

type ipinfoClient struct {
	v4Endpoint string
	v6Endpoint string
}

func (ipinfo ipinfoClient) getIpAddress(ipv6 bool) (string, error) {
	ipinfoURL := ipinfo.v4Endpoint
	if ipv6 {
		ipinfoURL = ipinfo.v6Endpoint
	}
	response, err := http.Get(ipinfoURL)
	if err != nil {
		return "", fmt.Errorf("error calling ipinfo endpoint %s: %w", ipinfoURL, err)
	}

	b, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body from ipinfo endpoint %s: %w", ipinfoURL, err)
	}

	data := IPInfoData{}
	err = json.Unmarshal(b, &data)
	if err != nil {
		return "", fmt.Errorf("error unmarshalling response from ipinfo endpoint %s: %w", ipinfoURL, err)
	}

	return data.IPAddress, nil
}

func getIpAddresses(c configuration, client IPInfoClient) (string, string) {
	ipv4needed, ipv6needed := ipNeeded(c)
	ipv4address := ""
	if ipv4needed {
		addr, err := client.getIpAddress(false)
		if err != nil {
			resolveErrorsTotal.WithLabelValues("A").Inc()
			log.Errorf("Unable to get ipv4 address: %v", err)
		}
		ipv4address = addr
		log.Debugf("Current host ipv4: %s", ipv4address)
	} else {
		log.Debug("No ipv4 needed")
	}
	ipv6address := ""
	if ipv6needed {
		addr, err := client.getIpAddress(true)
		if err != nil {
			resolveErrorsTotal.WithLabelValues("AAAA").Inc()
			log.Errorf("Unable to get ipv6 address: %v", err)
		}
		ipv6address = addr
		log.Debugf("Current host ipv6: %s", ipv6address)
	} else {
		log.Debug("No ipv6 needed")
	}
	return ipv4address, ipv6address
}
