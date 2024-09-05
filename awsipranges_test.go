package awsipranges

import (
	"encoding/json"
	"net"
	"os"
	"reflect"
	"testing"
)

// newFromFile reads the provided file and parses it
func newFromFile(f string) (*AWSIPRanges, error) {
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

func TestAWSIPRanges_Contains(t *testing.T) {
	tests := []struct {
		name    string
		ip      net.IP
		want    []Prefix
		wantErr bool
	}{
		{"empty", net.ParseIP(""), nil, true},
		{"non-aws", net.ParseIP("1.1.1.1"), nil, true},
		{"aws", net.ParseIP("3.5.12.4"),
			[]Prefix{
				{IPPrefix: "3.5.0.0/19", Region: "us-east-1", NetworkBorderGroup: "us-east-1", Service: "AMAZON"},
				{IPPrefix: "3.5.0.0/19", Region: "us-east-1", NetworkBorderGroup: "us-east-1", Service: "S3"},
				{IPPrefix: "3.5.0.0/19", Region: "us-east-1", NetworkBorderGroup: "us-east-1", Service: "EC2"},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a, err := newFromFile("testdata/ip-ranges-test.json")
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
