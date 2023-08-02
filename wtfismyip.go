package main

import (
	"encoding/json"
	"io"
	"net/http"
)

// WTFIsMyIPData is a data representation with the same structure returned by https://wtfismyip.com/json
type WTFIsMyIPData struct {
	YourFuckingIPAddress string `json:"YourFuckingIPAddress"`
}

func getIpAddress(ipv6 bool) (string, error) {
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
