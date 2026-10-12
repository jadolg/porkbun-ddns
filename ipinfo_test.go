package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type fakeIPInfoClient struct {
}

func (ipinfo fakeIPInfoClient) getIpAddress(ipv6 bool) (string, error) {
	if ipv6 {
		return "9afb:8b5f:078a:b914:e59e:f927:b0fa:ff31", nil
	}
	return "194.52.151.172", nil
}

func TestUnusedFromRecordsProduceEmptyIpAddresses(t *testing.T) {
	c := configuration{
		Records: []Record{
			{
				Domain: "example.com",
				Host:   "www",
				IpV6:   false,
				IpV4:   false,
			},
		},
	}
	ipv4, ipv6 := getIpAddresses(c, fakeIPInfoClient{})
	assert.Equal(t, "", ipv4)
	assert.Equal(t, "", ipv6)
}

func TestUsedFromRecordsProduceIPPaddresses(t *testing.T) {
	c := configuration{
		Records: []Record{
			{
				Domain: "example.com",
				Host:   "www",
				IpV6:   true,
				IpV4:   true,
			},
		},
	}
	ipv4, ipv6 := getIpAddresses(c, fakeIPInfoClient{})
	assert.Equal(t, "194.52.151.172", ipv4)
	assert.Equal(t, "9afb:8b5f:078a:b914:e59e:f927:b0fa:ff31", ipv6)
}
