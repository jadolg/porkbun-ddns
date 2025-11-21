# porkbun-ddns

[![CodeQL](https://github.com/jadolg/porkbun-ddns/actions/workflows/github-code-scanning/codeql/badge.svg)](https://github.com/jadolg/porkbun-ddns/actions/workflows/github-code-scanning/codeql) [![Go](https://github.com/jadolg/porkbun-ddns/actions/workflows/build.yml/badge.svg)](https://github.com/jadolg/porkbun-ddns/actions/workflows/build.yml) [![govulncheck](https://github.com/jadolg/porkbun-ddns/actions/workflows/govulncheck.yml/badge.svg)](https://github.com/jadolg/porkbun-ddns/actions/workflows/govulncheck.yml)

This project updates your DNS records at [Porkbun](https://porkbun.com/) with your current public IP address.
Your IP address is determined using https://wtfismyip.com/

Follow the white rabbit: https://dev.to/jadolg/keep-porkbun-dns-records-updated-automatically-with-your-current-ip-3539

## Install

### Snapcraft

[![Get it from the Snap Store](https://snapcraft.io/static/images/badges/en/snap-store-black.svg)](https://snapcraft.io/porkbun-ddns)

- Run `snap install porkbun-ddns`
- Place your configuration at `/var/snap/porkbun-ddns/common/config.yaml`
- An example configuration can be found at https://raw.githubusercontent.com/jadolg/porkbun-ddns/main/example.config.yaml
- Restart the service `snap restart porkbun-ddns.porkbun-ddns-daemon`

### Docker

Modify the following command to match your specific requirements (path to the config file):
```
docker run -d --restart always --network host -v $PWD/config.yaml:/config.yaml ghcr.io/jadolg/porkbun-ddns:latest
```
*note:* I'm using `--network host` to avoid setting up IPv6 in Docker.

### Debian/Ubuntu

```bash
wget -O - https://deb.akiel.dev/gpg.pub.key | gpg --dearmor | sudo tee /etc/apt/trusted.gpg.d/akiel.gpg > /dev/null
sudo apt-add-repository deb "https://deb.akiel.dev/ all main"
sudo apt install porkbun-ddns
```
or download the **.deb** from the Releases page and manually install with dpkg

- The config file is under /etc/porbun-ddns/config.yaml
- Remember to restart the service so it picks up the new configuration after updating it `sudo systemctl restart porbun-ddns`

### Direct download

Not a fan of running snaps or docker? Not a problem!
Grab the binary for your system/architecture directly from the [release](https://github.com/jadolg/porkbun-ddns/releases) page and start using it.

### From source

You can also just install from GitHub using Go directly `go install github.com/jadolg/porkbun-ddns@latest`

## Metrics

To enable the collection of metrics, change the config file and set a valid port.
The current metrics are:
- porkbun_update_errors_total
- porkbun_credentials_errors_total
- porkbun_connection_errors_total
- porkbun_resolve_errors_total
- porkbun_update_success_total
