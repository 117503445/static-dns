package dns

import (
	"fmt"
	"net"
	"path/filepath"

	"github.com/117503445/static-dns/pkg/cli"
	"github.com/miekg/dns"
	"github.com/rs/zerolog/log"
)

type DNSServer struct {
	udpServer *dns.Server
}

// HandleOutbound performs a DNS lookup for the provided domain and returns the first A record found.
func HandleOutbound(domain string) (dest string, err error) {
	// 创建一个新的 DNS 查询消息
	msg := new(dns.Msg)
	msg.SetQuestion(dns.Fqdn(domain), dns.TypeA)
	// 设置为递归查询（RD=1）
	msg.RecursionDesired = true

	// 创建客户端
	client := new(dns.Client)

	// 发送请求到本地 DNS 服务器
	resp, _, err := client.Exchange(msg, cli.Cli.Upstream)
	if err != nil {
		return "", err
	}

	// 解析响应中的答案部分
	for _, answer := range resp.Answer {
		switch rr := answer.(type) {
		case *dns.A:
			return net.IP(rr.A).String(), nil
		default:
			// 忽略其他类型的记录
		}
	}

	// 如果没有找到任何 A 记录，则返回空字符串
	return "", nil
}

func HandleStatic(domain string) (dest string) {
	for _, rule := range cli.Cli.Rules {
		matched, err := filepath.Match(rule.Pattern, domain)
		if err != nil {
			log.Error().Err(err).Msg("Failed to match pattern")
		} else {
			if matched {
				dest = rule.Dest
				break
			}
		}
	}
	return
}

func NewServer() *DNSServer {
	dns.HandleFunc(".", func(w dns.ResponseWriter, r *dns.Msg) {
		resp := new(dns.Msg)
		resp.SetReply(r)
		resp.Authoritative = true

		names := []string{}

		// TODO: Handle _acme-challenge
		for _, q := range resp.Question {
			names = append(names, q.Name)
			if q.Qtype == dns.TypeA {
				var err error
				dest := HandleStatic(q.Name)
				if dest != "" {
					log.Info().Str("name", q.Name).Str("dest", dest).Msg("by static")
				} else {
					dest, err = HandleOutbound(q.Name)
					if err != nil || dest != "" {
						log.Warn().Err(err).Msg("failed to get by outbound")
					} else {
						log.Info().Str("name", q.Name).Str("dest", dest).Msg("by outbound")
					}
				}

				if dest != "" {
					a := &dns.A{
						Hdr: dns.RR_Header{Name: q.Name, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 60},
						A:   net.ParseIP(dest),
					}
					resp.Answer = append(resp.Answer, a)
				}
			}
		}
		if len(resp.Answer) == 0 {
			log.Warn().Strs("names", names).Msg("No answer")
			resp.SetRcode(r, dns.RcodeNameError)
		}
		w.WriteMsg(resp)
	})

	return &DNSServer{}
}

func (s *DNSServer) Start() {
	port := cli.Cli.Port
	addr := fmt.Sprintf(":%d", port)
	udpServer := &dns.Server{Addr: addr, Net: "udp"}
	s.udpServer = udpServer

	log.Info().Int("port", port).Msg("Starting dns server")
	if err := s.udpServer.ListenAndServe(); err != nil {
		log.Fatal().Err(err).Msg("Failed to set udp listener")
	}
}
