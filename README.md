# porkbun-ddns

This project updates your DNS records at Porkbun with your current public IP address.
Your IP address is determined using https://wtfismyip.com/

## Install

- Run `snap install porkbun-ddns`
- Place your configuration at `/var/snap/porkbun-ddns/common/config.yaml`
- An example configuration can be found at https://raw.githubusercontent.com/jadolg/porkbun-ddns/main/example.config.yaml
- Restart the service `snap restart porkbun-ddns.porkbun-ddns-daemon`

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