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

	if flag.NArg() > 1 {
		log.Fatal("unexpected number of args")
	}

	var filters []awsipranges.Filter

	if ip := flag.Arg(0); ip != "" {
		filters = append(filters, awsipranges.Filter{
			Type:  awsipranges.FilterTypeIP,
			Value: ip,
		})
	}

	if region != "" {
		filters = append(filters, awsipranges.Filter{
			Type:  awsipranges.FilterTypeRegion,
			Value: region,
		})
	}

	if service != "" {
		filters = append(filters, awsipranges.Filter{
			Type:  awsipranges.FilterTypeService,
			Value: service,
		})
	}

	if networkBorderGroup != "" {
		filters = append(filters, awsipranges.Filter{
			Type:  awsipranges.FilterTypeNetworkBorderGroup,
			Value: networkBorderGroup,
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
