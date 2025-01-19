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

func NewServer() *DNSServer {
	dns.HandleFunc(".", func(w dns.ResponseWriter, r *dns.Msg) {
		// log.Info().Msg("Got request")
		resp := new(dns.Msg)
		resp.SetReply(r)
		resp.Authoritative = true

		names := []string{}

		for _, q := range resp.Question {
			names = append(names, q.Name)
			if q.Qtype == dns.TypeA {
				dest := ""

				for _, rule := range cli.Cli.Rules {
					matched, err := filepath.Match(rule.Pattern, q.Name)
					if err != nil {
						log.Error().Err(err).Msg("Failed to match pattern")
					} else {
						if matched {
							dest = rule.Dest
						}
					}
				}
				log.Info().Str("name", q.Name).Str("dest", dest).Msg("Match rule")
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
