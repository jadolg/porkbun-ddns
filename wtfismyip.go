package main

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
)

// WTFIsMyIPData is a data representation with the same structure returned by https://wtfismyip.com/json
type WTFIsMyIPData struct {
	YourFuckingIPAddress string `json:"YourFuckingIPAddress"`
}

type WTFIsMyIPClient interface {
	getIpAddress(ipv6 bool) (string, error)
}

type wtfismyipClient struct {
}

func (wtfismyip wtfismyipClient) getIpAddress(ipv6 bool) (string, error) {
	wtfismyipURL := "https://ipv4.wtfismyip.com/json"
	if ipv6 {
		wtfismyipURL = "https://ipv6.wtfismyip.com/json"
	}
	response, err := http.Get(wtfismyipURL)
	if err != nil {
		return "", err
	}

	b, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	data := WTFIsMyIPData{}
	err = json.Unmarshal(b, &data)
	if err != nil {
		return "", err
	}

	return data.YourFuckingIPAddress, nil
}

func getIpAddresses(c configuration, wtfismyip WTFIsMyIPClient) (string, string) {
	ipv4needed, ipv6needed := ipNeeded(c)
	ipv4address := ""
	if ipv4needed {
		wtfIPv4address, err := wtfismyip.getIpAddress(false)
		if err != nil {
			log.Error(err)
		}
		ipv4address = wtfIPv4address
		log.Debugf("Current host ipv4: %s", ipv4address)
	} else {
		log.Debug("No ipv4 needed")
	}
	ipv6address := ""
	if ipv6needed {
		wtfIPv6address, err := wtfismyip.getIpAddress(true)
		if err != nil {
			log.Error(err)
		}
		ipv6address = wtfIPv6address
		log.Debugf("Current host ipv6: %s", ipv6address)
	} else {
		log.Debug("No ipv6 needed")
	}
	return ipv4address, ipv6address
}
