# awsipranges

Helpers for working with public AWS IP range data.

Inspired by a much better existing version of this tool, [cmlccie/awsipranges](https://github.com/cmlccie/awsipranges), but written as a Go library for some eventual use cases that would require it.

## CLI

Originally this was only going to be a library, but a basic main function was helpful for testing and morphed into a simple CLI.

### Installation

Via `go install`:

```console
go install github.com/jar-b/awsipranges/cmd/awsipranges@latest
```

### Usage

```console
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

The output of a filtered query is printed as a table:

```console
% awsipranges 52.119.252.5
 |IP Prefix       |Region    |Network Border Group |Service  |
 |---------       |------    |-------------------- |-------  |
 |52.119.252.0/22 |us-west-2 |us-west-2            |AMAZON   |
 |52.119.252.0/22 |us-west-2 |us-west-2            |DYNAMODB |
```

### Examples

Search for the range of a specific IP address:

```console
awsipranges 52.119.252.5
```

List IP ranges for a region:

```console
awsipranges -region=us-west-2
```

List IP ranges for a service within a region:

```console
awsipranges -region=us-west-2 -service=DYNAMODB
```

Set the cachefile to expire after `240` hours (10 days):

```console
awsipranges -expiration=240h 52.119.252.5
```

## Library

`import github.com/jar-b/awsipranges`

### Usage

_TODO_
