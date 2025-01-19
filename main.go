package main

import (
	"strings"

	"github.com/117503445/goutils"
	"github.com/117503445/static-dns/pkg/cli"
	"github.com/117503445/static-dns/pkg/dns"
	"github.com/alecthomas/kong"
	kongtoml "github.com/alecthomas/kong-toml"
	"github.com/rs/zerolog/log"
)

func main() {
	goutils.InitZeroLog()

	kong.Parse(&cli.Cli, kong.Configuration(kongtoml.Loader, "./config.toml"))

	for _, rule := range cli.Cli.Rules {
		if rule.Type == "" {
			rule.Type = "glob"
		}
		if rule.Type == "glob" {
			if !strings.HasSuffix(rule.Pattern, "."){
				rule.Pattern += "?"
			}
		}
	}
	log.Info().Interface("cli", cli.Cli).Msg("")

	dns.NewServer().Start()
}
