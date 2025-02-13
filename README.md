# static-dns

> A static DNS resolution server

- If a DNS request matches the rules, it returns the specified IP address.
- If a DNS request does not match the rules, it forwards the request to an upstream DNS server.

## Quick Start

Prepare the configuration file and Docker Compose declaration file.

```sh
git clone https://github.com/117503445/static-dns.git
cd static-dns/docs/example
```

Run the server.

```sh
docker-compose up -d
```

Query DNS.

```sh
dig @127.0.0.1 -p 5053 router.lan
# router.lan. 0 IN A 192.168.1.1

dig @127.0.0.1 -p 5053 baidu.com
# baidu.com. 0 IN A 39.156.66.10
```

## Configuration Reference

The configuration file uses TOML format. Here is an example based on the quick start configuration:

```toml
# port = 5053
# upstream = "223.5.5.5:53"

[[rules]]
pattern = "*router.lan"
dest = "192.168.1.1"
```

`port` specifies the port for the DNS server, defaulting to 5053.

`upstream` specifies the address of the upstream DNS server, defaulting to 223.5.5.5:53.

`rules` specify the rules for the DNS server. Each rule contains the following fields:

- `pattern` is used for matching domain names in DNS requests.
- `type` indicates the type of matching rule, with the default being `glob`, and currently only `glob` is supported.
- `dest` specifies the IP address returned by the DNS server for matched requests.
