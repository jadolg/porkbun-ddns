# porkbun-ddns

This project updates your DNS records at Porkbun with your current public IP address.
Your IP address is determined using https://wtfismyip.com/

## Install

- Run `snap install porkbun-ddns`
- Place your configuration at `/var/snap/porkbun-ddns/common/config.yaml`
- Restart the service `snap restart porkbun-ddns`

## Example configuration

```yaml
porkbun_api_key: changeme
porkbun_secret_key: changeme

update_interval_minutes: 5

records:
  - domain: example.com
    host: www
    ipv6: true
    ipv4: true
  - domain: example.com
    host: blog
    ipv4: true
```