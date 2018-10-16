package main

import (
	"strings"
	"testing"

	"github.com/miekg/dns"
)

func TestLookup(t *testing.T) {

    options.Name = "e.root-servers.net"
	options.Dnssec = false
	options.Verbose = false

	conf_string := strings.NewReader("nameserver 8.8.8.8")
	config, _ := dns.ClientConfigFromReader(conf_string)

	ms, err := Lookup(config)
	if err != nil {
		t.Errorf(err.Error())
	}
    t.Logf("Found %s in %dms", options.Name, ms)

}
