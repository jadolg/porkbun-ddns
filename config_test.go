package main

import (
	"testing"
)

func TestPorkbunCredentials_invalid(t *testing.T) {
	tests := []struct {
		name        string
		credentials PorkbunCredentials
		want        bool
	}{
		{
			"When API Key is empty then return true",
			PorkbunCredentials{
				PorkbunAPIKey:    "",
				PorkbunSecretKey: "AAA",
			},
			true,
		},
		{
			"When Secret Key is empty then return true",
			PorkbunCredentials{
				PorkbunAPIKey:    "AAA",
				PorkbunSecretKey: "",
			},
			true,
		},
		{
			"When all data is passed correctly then return false",
			PorkbunCredentials{
				PorkbunAPIKey:    "AAA",
				PorkbunSecretKey: "AAAA",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := PorkbunCredentials{
				PorkbunAPIKey:    tt.credentials.PorkbunAPIKey,
				PorkbunSecretKey: tt.credentials.PorkbunSecretKey,
			}
			if got := c.invalid(); got != tt.want {
				t.Errorf("invalid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRecord_invalid(t *testing.T) {
	tests := []struct {
		name   string
		record Record
		want   bool
	}{
		{
			"When domain is empty then return true",
			Record{
				Domain: "",
			},
			true,
		},
		{
			"When host is empty then return true",
			Record{
				Host: "",
			},
			true,
		},
		{
			"When credentials is empty then return true",
			Record{
				Credentials: "",
			},
			true,
		},
		{
			"When both ipv4 and ipv6 are falsethen return true",
			Record{
				Host:   "AAA",
				Domain: "AAA",
				IpV6:   false,
				IpV4:   false,
			},
			true,
		},
		{
			"When all data is passed correctly then return false",
			Record{
				Host:        "AAA",
				Domain:      "AAA",
				IpV6:        true,
				IpV4:        true,
				Credentials: "AA",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Record{
				Domain:      tt.record.Domain,
				Host:        tt.record.Host,
				IpV6:        tt.record.IpV6,
				IpV4:        tt.record.IpV4,
				Credentials: tt.record.Credentials,
			}
			if got := r.invalid(); got != tt.want {
				t.Errorf("invalid() = %v, want %v", got, tt.want)
			}
		})
	}
}
