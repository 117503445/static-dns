package cli

type Rule struct {
	Type    string
	Pattern string
	Dest    string
}

var Cli struct {
	Port     int     `default:"5053" help:"Port to listen."`
	Upstream string  `default:"223.5.5.5:53" help:"Upstream DNS server"`
	Rules    []*Rule `help:"Rules to match."`
}
