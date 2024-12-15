# awsipranges

Helpers for working with public AWS IP range data.

## CLI

### Installation

Via `go install`:

```sh
go install github.com/jar-b/awsipranges/cmd/awsipranges@latest
```

### Usage

```
% awsipranges -h
Check whether an IP address is in an AWS range.

Usage: awsipranges [flags] [ip]

Flags:
  -cachefile string
        Location of the cached ip-ranges.json file (default "~/.aws/ip-ranges.json")
  -expiration string
        Duration after which the cached ranges file should be replaced
  -network-border-group string
        Network border group to filter on (e.g. us-west-2-lax-1)
  -region string
        Region name to filter on (e.g. us-east-1)
  -service string
        Service name to filter on (e.g. EC2)
```

### Examples

Check for a specific IP:

```sh
awsipranges 52.119.252.5
```

List IP ranges for a region:

```sh
awsipranges -region=us-west-2
```

List IP ranges for a service within a region:

```sh
awsipranges -region=us-west-2 -service=DYNAMODB
```

## Library

`import github.com/jar-b/awsipranges`

### Usage

_TODO_
