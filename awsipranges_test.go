package awsipranges

import (
	"encoding/json"
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

func TestAWSIPRanges_Filter(t *testing.T) {
	tests := []struct {
		name    string
		filters []Filter
		want    []Prefix
		wantErr bool
	}{
		{
			name: "invalid",
			filters: []Filter{
				{
					Type:  "invalid",
					Value: "",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "ip",
			filters: []Filter{
				{
					Type:  FilterTypeIP,
					Value: "52.94.76.1",
				},
			},
			want: []Prefix{
				{
					IPPrefix:           "52.94.76.0/22",
					Region:             "us-west-2",
					Service:            "AMAZON",
					NetworkBorderGroup: "us-west-2",
				},
			},
		},
		{
			name: "network border group",
			filters: []Filter{
				{
					Type:  FilterTypeNetworkBorderGroup,
					Value: "us-west-2",
				},
			},
			want: []Prefix{
				{
					IPPrefix:           "52.94.76.0/22",
					Region:             "us-west-2",
					Service:            "AMAZON",
					NetworkBorderGroup: "us-west-2",
				},
			},
		},
		{
			name: "region",
			filters: []Filter{
				{
					Type:  FilterTypeRegion,
					Value: "us-west-2",
				},
			},
			want: []Prefix{
				{
					IPPrefix:           "52.94.76.0/22",
					Region:             "us-west-2",
					Service:            "AMAZON",
					NetworkBorderGroup: "us-west-2",
				},
			},
		},
		{
			name: "service",
			filters: []Filter{
				{
					Type:  FilterTypeService,
					Value: "CODEBUILD",
				},
			},
			want: []Prefix{
				{
					IPPrefix:           "3.101.177.48/29",
					Region:             "us-west-1",
					Service:            "CODEBUILD",
					NetworkBorderGroup: "us-west-1",
				},
			},
		},
		{
			name: "all",
			filters: []Filter{
				{
					Type:  FilterTypeIP,
					Value: "3.101.177.48",
				},
				{
					Type:  FilterTypeService,
					Value: "CODEBUILD",
				},
				{
					Type:  FilterTypeRegion,
					Value: "us-west-1",
				},
				{
					Type:  FilterTypeNetworkBorderGroup,
					Value: "us-west-1",
				},
			},
			want: []Prefix{
				{
					IPPrefix:           "3.101.177.48/29",
					Region:             "us-west-1",
					Service:            "CODEBUILD",
					NetworkBorderGroup: "us-west-1",
				},
			},
		},
		{
			name: "no match",
			filters: []Filter{
				{
					Type:  FilterTypeRegion,
					Value: "us-west-99",
				},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a, err := newFromFile("testdata/ip-ranges-test.json")
			if err != nil {
				t.Fatalf("reading testdata: %v", err)
			}

			got, err := a.Filter(tt.filters)
			if (err != nil) != tt.wantErr {
				t.Errorf("AWSIPRanges.Filter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AWSIPRanges.Filter() = %v, want %v", got, tt.want)
			}
		})
	}
}
