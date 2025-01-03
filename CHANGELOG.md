# Changelog

## v0.0.6

### Fixed
- Properly omit items when one filter matches but another does not

## v0.0.5

### Changed
- `cli`: IP address is now passed as a flag (`-ip`) for consistency with other filters

## v0.0.4

### Changed
- `Filter` now includes a `Values` field (type `[]string`), replacing `Value` (type `string`)

## v0.0.3

### Added
- Support case insensitive matching on network border group, region, and service

## v0.0.2

### Fixed
- Fixed installation from the `cmd/awsipranges` subdirectory

## v0.0.1
- Initial release
