# vigor
DNS resolution metrics for Graphite

Looks up a record with or without EDNS/DNSSEC and shove the response time into a graphite server or something that handles graphit formatted metrics.

```
Usage:
  vigor [OPTIONS]

Application Options:
  -v, --verbose   verbose output
  -s, --dnssec    request DNSSEC
  -c, --config=   Resolver config file (default /etc/resolv.conf) (default:
                  /etc/resolv.conf)
  -g, --graphite= Graphite host and port eg. 10.5.4.3:2003
  -p, --path=     Graphite Metric Path
  -n, --name=

Help Options:
  -h, --help      Show this help message
  ```

# TODO

- Allow for multiple lookups, multiple hostnames
- Move name to a proper argument vs. option
