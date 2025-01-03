// Package awsipranges provides helpers for working with public AWS IP range data
package awsipranges

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"slices"
	"strings"
)

const ipRangesURL = "https://ip-ranges.amazonaws.com/ip-ranges.json"

// AWSIPRanges stores the content from `ip-ranges.json`
//
// See the AWS documentation for additional details:
// - https://docs.aws.amazon.com/vpc/latest/userguide/aws-ip-work-with.html
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

// FilterType is the type of filter to apply while iterating over IP range
// data
//
// Filtering can be done on IP address, network border group, region, and
// service.
type FilterType string

const (
	FilterTypeIP                 FilterType = "ip"
	FilterTypeNetworkBorderGroup FilterType = "network-border-group"
	FilterTypeRegion             FilterType = "region"
	FilterTypeService            FilterType = "service"
)

// Filter stores the filter type and values used to filter results
type Filter struct {
	Type   FilterType
	Values []string
}

// Get returns the content from the latest "ip-ranges.json" file
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
				keep = slices.ContainsFunc(f.Values, func(e string) bool {
					ip := net.ParseIP(e)
					_, ipNet, err := net.ParseCIDR(p.IPPrefix)
					if err != nil {
						// if the IP prefix cannot be parsed, proceed without filtering
						return keep
					}

					return ipNet.Contains(ip)
				})
			case FilterTypeNetworkBorderGroup:
				keep = slices.ContainsFunc(f.Values, func(e string) bool {
					return strings.EqualFold(e, p.NetworkBorderGroup)
				})
			case FilterTypeRegion:
				keep = slices.ContainsFunc(f.Values, func(e string) bool {
					return strings.EqualFold(e, p.Region)
				})
			case FilterTypeService:
				keep = slices.ContainsFunc(f.Values, func(e string) bool {
					return strings.EqualFold(e, p.Service)
				})
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
