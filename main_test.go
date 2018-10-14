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


	expected := "192.203.230.10"

	conf_string := strings.NewReader("nameserver 8.8.8.8")
	config, _ := dns.ClientConfigFromReader(conf_string)

	rec, ms := Lookup(config)
	if !strings.HasSuffix(rec, expected) {
		t.Errorf("Lookup failed:\n\tfound:%s\n\texpected:%s\n", rec, expected)
	}
    t.Logf("Found %s in %dms", options.Name, ms)

}
