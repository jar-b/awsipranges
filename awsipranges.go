package awsipranges

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
)

const ipRangesURL = "https://ip-ranges.amazonaws.com/ip-ranges.json"

// AWSIPRanges stores the content from `ip-ranges.json`
//
// Ref: https://docs.aws.amazon.com/vpc/latest/userguide/aws-ip-work-with.html
type AWSIPRanges struct {
	SyncToken    string       `json:"syncToken"`
	CreateDate   string       `json:"createDate"`
	Prefixes     []Prefix     `json:"prefixes"`
	IPV6Prefixes []IPV6Prefix `json:"ipv6_prefixes"`
}

type Prefix struct {
	IPNet              net.IPNet
	IPPrefix           string `json:"ip_prefix"`
	Region             string `json:"region"`
	NetworkBorderGroup string `json:"network_border_group"`
	Service            string `json:"service"`
}

type IPV6Prefix struct {
	IPV6Net            net.IPNet
	IPV6Prefix         string `json:"ipv6_prefix"`
	Region             string `json:"region"`
	NetworkBorderGroup string `json:"network_border_group"`
	Service            string `json:"service"`
}

// NewFromFile reads the provided file and parses it
func NewFromFile(f string) (*AWSIPRanges, error) {
	b, err := os.ReadFile(f)
	if err != nil {
		return nil, err
	}

	var ranges AWSIPRanges
	if err := json.Unmarshal(b, &ranges); err != nil {
		return nil, err
	}

	return &ranges, nil
}

// New fetches the latest `ip-ranges.json` file and parses it
func New() (*AWSIPRanges, error) {
	resp, err := http.Get(ipRangesURL)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var ranges AWSIPRanges
	if err := json.Unmarshal(b, &ranges); err != nil {
		return nil, err
	}

	return &ranges, nil
}

func (a *AWSIPRanges) Contains(ip net.IP) (*Prefix, error) {
	for _, p := range a.Prefixes {
		if p.IPNet.Contains(ip) {
			return &p, nil
		}
	}

	return nil, NewNotInRangeError(ip)
}

type NotInRangeError struct {
	ip net.IP
}

func (e NotInRangeError) Error() string {
	return fmt.Sprintf("%s not in AWS IP ranges", e.ip.String())
}

func NewNotInRangeError(ip net.IP) error {
	return NotInRangeError{ip: ip}
}
