package awsipranges

import (
	"net"
	"reflect"
	"testing"
)

func TestAWSIPRanges_Contains(t *testing.T) {
	tests := []struct {
		name    string
		ip      net.IP
		want    *Prefix
		wantErr bool
	}{
		{"empty", net.ParseIP(""), nil, true},
		{"non-aws", net.ParseIP("1.1.1.1"), nil, true},
		{"aws", net.ParseIP("3.4.12.4"),
			&Prefix{
				IPPrefix:           "3.4.12.4/32",
				Region:             "eu-west-1",
				Service:            "AMAZON",
				NetworkBorderGroup: "eu-west-1",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a, err := NewFromFile("testdata/ip-ranges-20240902.json")
			if err != nil {
				t.Fatalf("reading testdata: %v", err)
			}

			got, err := a.Contains(tt.ip)
			if (err != nil) != tt.wantErr {
				t.Errorf("AWSIPRanges.Contains() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AWSIPRanges.Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}
