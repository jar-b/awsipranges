package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"text/tabwriter"
	"time"

	"github.com/jar-b/awsipranges"
)

const (
	cachefilePath = ".aws/ip-ranges.json"

	// createDateFormat is the format of the `createDate` field in the
	// underlying JSON (YY-MM-DD-hh-mm-ss)
	//
	// Ref: https://docs.aws.amazon.com/vpc/latest/userguide/aws-ip-syntax.html
	createDateFormat = "2006-01-02-15-04-05"
)

var (
	cachefile          string
	ip                 string
	networkBorderGroup string
	region             string
	service            string
	expiration         string
)

func main() {
	// slightly better usage output
	flag.Usage = func() {
		fmt.Fprint(flag.CommandLine.Output(), "Check whether an IP address is in an AWS range.\n\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [flags]\n\nFlags:\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.StringVar(&cachefile, "cachefile", defaultCachefilePath(), "Location of the cached ip-ranges.json file")
	flag.StringVar(&expiration, "expiration", "", "Duration after which the cached ranges file should be replaced")
	flag.StringVar(&ip, "ip", "", "IP address to filter on (e.g. 1.2.3.4)")
	flag.StringVar(&networkBorderGroup, "network-border-group", "", "Network border group to filter on (e.g. us-west-2-lax-1)")
	flag.StringVar(&region, "region", "", "Region name to filter on (e.g. us-east-1)")
	flag.StringVar(&service, "service", "", "Service name to filter on (e.g. EC2)")
	flag.Parse()

	if flag.NArg() > 0 {
		log.Fatal("unexpected number of arguments")
	}

	var filters []awsipranges.Filter

	if ip != "" {
		filters = append(filters, awsipranges.Filter{
			Type:   awsipranges.FilterTypeIP,
			Values: []string{ip},
		})
	}

	if region != "" {
		filters = append(filters, awsipranges.Filter{
			Type:   awsipranges.FilterTypeRegion,
			Values: []string{region},
		})
	}

	if service != "" {
		filters = append(filters, awsipranges.Filter{
			Type:   awsipranges.FilterTypeService,
			Values: []string{service},
		})
	}

	if networkBorderGroup != "" {
		filters = append(filters, awsipranges.Filter{
			Type:   awsipranges.FilterTypeNetworkBorderGroup,
			Values: []string{networkBorderGroup},
		})
	}

	if len(filters) == 0 {
		log.Fatal("must provide an IP argument or set the -network-border-group, -region, or -service flag")
	}

	ranges, err := loadRanges()
	if err != nil {
		log.Fatal(err)
	}

	matches, err := ranges.Filter(filters)
	if err != nil {
		log.Fatal(err)
	}

	write(matches)
}

// defaultCachefilePath constructs a default path to the cachefile
func defaultCachefilePath() string {
	u, _ := user.Current()
	return filepath.Join(u.HomeDir, cachefilePath)
}

// isExpired checks whether the createDate of a cached ip-ranges.json
// file is older than the configured expiration duration
//
// If expiration is not set, always returns false.
func isExpired(createDate string) (bool, error) {
	if expiration == "" {
		return false, nil
	}

	created, err := time.Parse(createDateFormat, createDate)
	if err != nil {
		return false, err
	}
	expirationDuration, err := time.ParseDuration(expiration)
	if err != nil {
		return false, err
	}

	if expirationDuration > time.Since(created) {
		return false, nil
	}

	fmt.Println("Cache is expired, refreshing")
	return true, nil
}

// loadRanges attempts to read ip-ranges data from cache, falling back
// to fetching the source file if os.ReadFile fails or the creation date
// exceeds the configured cache expiration time
func loadRanges() (*awsipranges.AWSIPRanges, error) {
	if b, err := os.ReadFile(cachefile); err == nil {
		var ranges awsipranges.AWSIPRanges
		if err := json.Unmarshal(b, &ranges); err == nil {
			if exp, err := isExpired(ranges.CreateDate); err != nil {
				return nil, err
			} else if !exp {
				return &ranges, nil
			}
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

func write(matches []awsipranges.Prefix) {
	if len(matches) == 0 {
		fmt.Println("No matches found.")
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)
	fmt.Fprintln(w, "\tIP Prefix\tRegion\tNetwork Border Group\tService\t")
	fmt.Fprintln(w, "\t---------\t------\t--------------------\t-------\t")
	for _, m := range matches {
		fmt.Fprintf(w, "\t%s\t%s\t%s\t%s\t\n", m.IPPrefix, m.Region, m.NetworkBorderGroup, m.Service)
	}
	w.Flush()
}
