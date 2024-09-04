package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/jar-b/awsipranges"
)

var (
	cachefile string = "~/.aws/ip-ranges.json"
	service   string
	region    string
)

func main() {
	// slightly better usage output
	flag.Usage = func() {
		fmt.Fprint(flag.CommandLine.Output(), "Check whether an IP address is in an AWS range.\n\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [flags] [ip]\n\nFlags:\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.StringVar(&cachefile, "cachefile", cachefile, "Location of the cached ip-ranges.json file")
	flag.StringVar(&region, "region", "", "Region name to filter on (e.g. us-east-1)")
	flag.StringVar(&service, "service", "", "Service name to filter on (e.g. EC2)")
	flag.Parse()

	if flag.NArg() == 0 && region == "" && service == "" {
		log.Fatal("must provide an IP argument or set the -region or -service flag")
	}
	if flag.NArg() > 1 {
		log.Fatal("unexpected number of args")
	}

	if ip := flag.Arg(0); ip != "" {
		if err := contains(ip); err != nil {
			log.Fatal(err)
		}
	}
	// TODO filter by region
	// TODO filter by service
	// TODO filter by region and service
}

func contains(s string) error {
	// TODO: check cache
	ranges, err := awsipranges.New()
	if err != nil {
		log.Fatal(err)
	}

	match, err := ranges.Contains(net.ParseIP(s))
	if err != nil {
		log.Fatal(err)
	}

	b, _ := json.MarshalIndent(match, "", "  ")
	fmt.Printf("%s", string(b))
	return nil
}
