package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/user"
	"path/filepath"

	"github.com/jar-b/awsipranges"
)

const cachefilePath = ".aws/ip-ranges.json"

var (
	cachefile          string
	networkBorderGroup string
	region             string
	service            string
)

func main() {
	// slightly better usage output
	flag.Usage = func() {
		fmt.Fprint(flag.CommandLine.Output(), "Check whether an IP address is in an AWS range.\n\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [flags] [ip]\n\nFlags:\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.StringVar(&cachefile, "cachefile", defaultCachefilePath(), "Location of the cached ip-ranges.json file")
	flag.StringVar(&networkBorderGroup, "network-border-group", "", "Network border group to filter on (e.g. us-west-2-lax-1)")
	flag.StringVar(&region, "region", "", "Region name to filter on (e.g. us-east-1)")
	flag.StringVar(&service, "service", "", "Service name to filter on (e.g. EC2)")
	flag.Parse()

	if flag.NArg() == 0 && region == "" && service == "" && networkBorderGroup == "" {
		log.Fatal("must provide an IP argument or set the -network-border-group, -region, or -service flag")
	}
	if flag.NArg() > 1 {
		log.Fatal("unexpected number of args")
	}

	ranges, err := loadRanges()
	if err != nil {
		log.Fatal(err)
	}

	if ip := flag.Arg(0); ip != "" {
		if err := contains(ranges, ip); err != nil {
			log.Fatal(err)
		}
	}

	// TODO filter by region, service, network border group
}

// defaultCachefilePath constructs a default path to the cachefile
func defaultCachefilePath() string {
	u, _ := user.Current()
	return filepath.Join(u.HomeDir, cachefilePath)
}

// loadRanges attempts to read ip-ranges data from cache, falling back
// to fetching the source file if os.ReadFile fails or the creation date
// exceeds the configured cache expiration time
func loadRanges() (*awsipranges.AWSIPRanges, error) {
	if b, err := os.ReadFile(cachefile); err == nil {
		var ranges awsipranges.AWSIPRanges
		if err := json.Unmarshal(b, &ranges); err == nil {
			// TODO: check configured expiration time
			return &ranges, nil
		}
	}

	b, err := awsipranges.Get()
	if err != nil {
		return nil, err
	}

	if err := os.WriteFile(cachefile, b, 0644); err != nil {
		//TODO: debug logging when cache write fails
	}

	var ranges awsipranges.AWSIPRanges
	if err := json.Unmarshal(b, &ranges); err != nil {
		return nil, err
	}

	return &ranges, nil
}

func contains(ranges *awsipranges.AWSIPRanges, s string) error {
	match, err := ranges.Contains(net.ParseIP(s))
	if err != nil {
		log.Fatal(err)
	}

	b, _ := json.MarshalIndent(match, "", "  ")
	fmt.Printf("%s", string(b))
	return nil
}
