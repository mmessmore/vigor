package main

import (
	"fmt"
	"log"
	"math"
	"net"
	"os"
	"time"

	"github.com/jessevdk/go-flags"
	"github.com/miekg/dns"
	"github.com/mmessmore/vigor/graphite"
)

type Options struct {
	Verbose    bool   `short:"v" long:"verbose" description:"verbose output"`
	Dnssec     bool   `short:"s" long:"dnssec" description:"request DNSSEC"`
	ConfigFile string `short:"c" long:"config" default:"/etc/resolv.conf" description:"Resolver config file"`
	Graphite   string `short:"g" long:"graphite" required:"1" description:"Graphite host and port eg. 10.5.4.3:2003"`
    GPath      string `short:"p" long:"path" default:"" description:"Graphite Metric Path (default: vigor.hostname.first_resolver"`
	Name       string `short:"n" long:"name" required:"1"`
}

var options = Options{}

func main() {

	config := parse_args()

	_, elapsed := Lookup(config)
	ms := int(math.Round(elapsed.Seconds() * 1000))

	if options.Verbose {
		fmt.Printf("That took %d milliseconds", ms)
	}
	graphite.SendMetric(options.Graphite, options.GPath, ms)
}

func parse_args() *dns.ClientConfig {
	_, err := flags.Parse(&options)
	if err != nil {
		os.Exit(22)
	}
	config, _ := dns.ClientConfigFromFile(options.ConfigFile)

	if options.GPath == "" {
        hostname, _ := os.Hostname()
		options.GPath = fmt.Sprintf("vigor.%s.%s", hostname, config.Servers[0])
	}

	return config
}

func Lookup(config *dns.ClientConfig) (string, time.Duration) {
	c := new(dns.Client)
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(options.Name), dns.TypeA)
	m.RecursionDesired = true
	if options.Dnssec {
		m.SetEdns0(4096, true)
	}

	if options.Verbose {
		log.Println("Info: Looking up:", options.Name)
	}
	start := time.Now()
	r, _, err := c.Exchange(m, net.JoinHostPort(config.Servers[0], config.Port))
	elapsed := time.Since(start)
	if r == nil {
		log.Fatalf("Error: %s\n", err.Error())
	}
	if r.Rcode != dns.RcodeSuccess {
		log.Fatalf("Error: Invalid answer name %s \n", options.Name)
	}

	if options.Verbose {
		for _, a := range r.Answer {
			fmt.Println(a)
		}
	}
	return fmt.Sprintf("%v", r.Answer[0]), elapsed
}
