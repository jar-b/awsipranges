package awsipranges

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
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

// Prefix represents a single entry in the prefixes list
type Prefix struct {
	IPPrefix           string `json:"ip_prefix"`
	Region             string `json:"region"`
	NetworkBorderGroup string `json:"network_border_group"`
	Service            string `json:"service"`
}

// IPV6Prefix represents a single entry in the ipv6_prefixes list
type IPV6Prefix struct {
	IPV6Prefix         string `json:"ipv6_prefix"`
	Region             string `json:"region"`
	NetworkBorderGroup string `json:"network_border_group"`
	Service            string `json:"service"`
}

type FilterType string

const (
	FilterTypeIP                 FilterType = "ip"
	FilterTypeNetworkBorderGroup FilterType = "network-border-group"
	FilterTypeRegion             FilterType = "region"
	FilterTypeService            FilterType = "service"
)

type Filter struct {
	Type  FilterType
	Value string
}

// Get fetches the latest "ip-ranges.json" file
func Get() ([]byte, error) {
	resp, err := http.Get(ipRangesURL)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

// New fetches the latest "ip-ranges.json" file and parses it
func New() (*AWSIPRanges, error) {
	b, err := Get()
	if err != nil {
		return nil, err
	}

	var ranges AWSIPRanges
	if err := json.Unmarshal(b, &ranges); err != nil {
		return nil, err
	}

	return &ranges, nil
}

// Filter returns all prefix entries which match the provided filters
func (a *AWSIPRanges) Filter(filters []Filter) ([]Prefix, error) {
	var prefixes []Prefix

	for _, p := range a.Prefixes {
		keep := true
		for _, f := range filters {
			switch f.Type {
			case FilterTypeIP:
				ip := net.ParseIP(f.Value)
				_, ipNet, err := net.ParseCIDR(p.IPPrefix)
				if err != nil {
					// if the IP prefix cannot be parsed, proceed without filtering
					continue
				}

				if !ipNet.Contains(ip) {
					keep = false
				}
			case FilterTypeNetworkBorderGroup:
				if f.Value != p.NetworkBorderGroup {
					keep = false
				}
			case FilterTypeRegion:
				if f.Value != p.Region {
					keep = false
				}
			case FilterTypeService:
				if f.Value != p.Service {
					keep = false
				}
			default:
				return nil, fmt.Errorf("invalid filter type")
			}
		}
		if keep {
			prefixes = append(prefixes, p)
		}
	}

	return prefixes, nil
}
