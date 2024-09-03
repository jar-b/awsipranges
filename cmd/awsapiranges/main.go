package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/jar-b/awsipranges"
)

var (
	cachefile string = "~/.aws/ip-ranges.json"
	ip        string

	downloadCmd = flag.NewFlagSet("download", flag.ExitOnError)
	containsCmd = flag.NewFlagSet("contains", flag.ExitOnError)
)

var subcommands = map[string]*flag.FlagSet{
	downloadCmd.Name(): downloadCmd,
	containsCmd.Name(): containsCmd,
}

// setup flags and usage specific to individual subcommands
func setupSubcommands() {
	containsCmd.StringVar(&ip, "ip", "", "IP address to check against AWS ranges")
	containsCmd.Usage = func() {
		fmt.Fprint(flag.CommandLine.Output(), "Check whether an IP address is in an AWS range.\n\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [flags]\n\nFlags:\n", os.Args[0])
		containsCmd.PrintDefaults()
	}

	downloadCmd.Usage = func() {
		fmt.Fprint(flag.CommandLine.Output(), "Download the ip-ranges.json file.\n\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [flags]\n\nFlags:\n", os.Args[0])
		downloadCmd.PrintDefaults()
	}
}

func setupGlobalFlags() {
	for _, cmd := range subcommands {
		cmd.StringVar(&cachefile, "cachefile", cachefile, "Location of the cached ip-ranges.json file")
	}
}
func main() {
	setupGlobalFlags()
	setupSubcommands()

	if len(os.Args) == 1 {
		log.Fatal("missing subcommand")
	}

	arg1 := os.Args[1]
	cmd := subcommands[arg1]
	if cmd == nil {
		log.Fatalf("unknown subcommand '%s'", os.Args[1])
	}

	if len(os.Args) > 2 {
		cmd.Parse(os.Args[2:])
	}

	switch arg1 {
	case downloadCmd.Name():
		if err := download(); err != nil {
			log.Fatalf("download: %w", err)
		}
	case containsCmd.Name():
		if err := contains(ip); err != nil {
			log.Fatalf("download: %w", err)
		}
	default:
		log.Fatalf("subcommand not implemented '%s'", arg1)
	}
}

func download() error {
	// TODO: implement
	return nil
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

	fmt.Printf("%+v\n", match)
	return nil
}
