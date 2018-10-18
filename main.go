package main

import (
	"errors"
	"fmt"
	"log"
	"math"
	"net"
	"os"
	"time"

	"github.com/jessevdk/go-flags"
	"github.com/miekg/dns"
)

type Options struct {
	Verbose        bool   `short:"v" long:"verbose" description:"verbose output"`
	Dnssec         bool   `short:"s" long:"dnssec" description:"request DNSSEC"`
	ConfigFile     string `short:"c" long:"config" default:"/etc/resolv.conf" description:"Resolver config file"`
	Graphite       string `short:"g" long:"graphite" required:"1" description:"Graphite host and port eg. 10.5.4.3:2003"`
	GPath          string `short:"p" long:"path" default:"" description:"Graphite Metric Path (default: vigor.hostname.first_resolver"`
	Name           string `short:"n" long:"name" required:"1"`
	ReportInterval int64  `short:"r" long:"report-interval" default:"10" description:"Report every x s"`
	QueryInterval  int64  `short:"q" long:"query-interval" default:"1" description:"Query every y s"`
}

type Work struct {
	Duration int
	Error    error
}

var options = Options{}
var WorkQueue = make(chan Work, 10000)

func main() {

	config := parse_args()

	go worker()
	collect(config)
}

func worker() {
	conn := GetGraphiteClient(options.Graphite)
	for {
		time.Sleep(time.Duration(options.ReportInterval) * time.Second)
		total := 0
		num := 0
		errs := 0
		high := 0
		low := 0
	Inner:
		for {
			select {
			case val := <-WorkQueue:
				total += val.Duration
				if val.Error != nil {
					errs++
				}
				num++
				if val.Duration > high {
					high = val.Duration
				}
				if val.Duration < low {
					low = val.Duration
				}
			default:
				if options.Verbose {
					fmt.Println("Sending metrics")
				}
				SendMetric(
					conn,
					fmt.Sprintf("%s.avg_ms", options.GPath),
					int(math.Round(float64(total)/float64(num))))
				SendMetric(
					conn,
					fmt.Sprintf("%s.total_ms", options.GPath),
					total)
				SendMetric(
					conn,
					fmt.Sprintf("%s.num", options.GPath),
					num)
				SendMetric(
					conn,
					fmt.Sprintf("%s.errors", options.GPath),
					errs)
				SendMetric(
					conn,
					fmt.Sprintf("%s.high", options.GPath),
					high)
				SendMetric(
					conn,
					fmt.Sprintf("%s.low", options.GPath),
					low)

				break Inner
			}
		}
	}
}

func collect(config *dns.ClientConfig) {
	for {
		time.Sleep(time.Duration(options.QueryInterval) * time.Second)
		elapsed, err := Lookup(config)
		ms := int(math.Round(elapsed.Seconds() * 1000))

		if options.Verbose {
			log.Printf("That took %d milliseconds\n", ms)
		}
		WorkQueue <- Work{ms, err}
	}
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

func Lookup(config *dns.ClientConfig) (time.Duration, error) {
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
		log.Printf("Error: %s\n", err.Error())
		return elapsed, err
	}
	if r.Rcode != dns.RcodeSuccess {
		message := fmt.Sprintf("Error: Invalid answer name %s", options.Name)
		log.Println(message)
		return elapsed, errors.New(message)
	}

	if options.Verbose {
		for _, a := range r.Answer {
			fmt.Println(a)
		}
	}
	return elapsed, nil
}

func GetGraphiteClient(dest string) net.Conn {
	conn, err := net.Dial("tcp", dest)
	if err != nil {
		log.Fatalln(err.Error())
	}
	return conn
}

func SendMetric(conn net.Conn, path string, metric int) {
	_, err := fmt.Fprintln(conn, path, metric, time.Now().Unix())
	if err != nil {
		log.Printf("Graphite Connection misdapeared.  Retrying once.")
		log.Printf("Once.")
		conn = GetGraphiteClient(options.Graphite)
		_, inerr := fmt.Fprintln(conn, path, metric, time.Now().Unix())
		if inerr != nil {
			log.Printf("Graphite's still gone.  Giving up hope.")
			log.Fatal(inerr.Error())
		}
	}
	fmt.Println(path, metric, time.Now().Unix())
}
