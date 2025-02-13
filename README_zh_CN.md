# static-dns

> 静态 DNS 解析服务器

- DNS 请求和规则匹配，则返回指定 IP 地址
- DNS 请求和规则不匹配，则向上游 DNS 服务器请求

## 快速开始

准备配置文件和 Docker Compose 声明文件

```sh
git clone https://github.com/117503445/static-dns.git
cd static-dns/docs/example
```

运行

```sh
docker-compose up -d
```

查询 DNS

```sh
dig @127.0.0.1 -p 5053 router.lan
# router.lan.	0	IN	A	192.168.1.1

dig @127.0.0.1 -p 5053 baidu.com
# baidu.com.	0	IN	A	39.156.66.10
```

## 配置参考

使用 TOML 格式配置文件，以快速开始中的配置为例

```toml
# port = 5053
# upstream = "223.5.5.5:53"

[[rules]]
pattern = "*router.lan"
dest = "192.168.1.1"
```

`port` 用于指定 DNS 服务器端口，默认为 5053

`upstream` 用于指定 DNS 服务器地址，默认为 223.5.5.5:53

`rules` 用于指定 DNS 服务器规则。每个 rule 都包含以下字段

- `pattern` 用于 DNS 请求的域名匹配规则
- `type` 匹配规则类型，默认为 `glob`，目前只支持 `glob`
- `dest` 用于指定 DNS 服务器返回的 IP 地址
