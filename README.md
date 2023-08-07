# porkbun-ddns

This project updates your DNS records at Porkbun with your current public IP address.
Your IP address is determined using https://wtfismyip.com/

## Install

### Snapcraft

[![Get it from the Snap Store](https://snapcraft.io/static/images/badges/en/snap-store-black.svg)](https://snapcraft.io/porkbun-ddns)

- Run `snap install porkbun-ddns`
- Place your configuration at `/var/snap/porkbun-ddns/common/config.yaml`
- An example configuration can be found at https://raw.githubusercontent.com/jadolg/porkbun-ddns/main/example.config.yaml
- Restart the service `snap restart porkbun-ddns.porkbun-ddns-daemon`

### Docker

Modify the following command to match your speciffic requirements (path to the config file):
```
docker run -d --restart always --network host -v $PWD/config.yaml:/config.yaml ghcr.io/jadolg/porkbun-ddns:latest
```
*note:* I'm using `--network host` to avoid setting up IPv6 in Docker.

## Example configuration

```yaml
update_interval_minutes: 5

credentials:
  prod:
    porkbun_api_key: changeme
    porkbun_secret_key: changeme
  dev:
    porkbun_api_key: changeme
    porkbun_secret_key: changeme

records:
  - domain: example.com
    host: www
    ipv6: true
    ipv4: true
    credentials: prod
  - domain: example.com
    host: blog
    ipv4: true
    credentials: dev

```
