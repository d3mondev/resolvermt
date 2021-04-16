# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]
- No changes yet.

## [0.3.0] - 2021-04-16
### Changed
- Optimize the number of goroutines created during concurrent Resolve calls.

## [0.2.0] - 2021-04-07
### Added
- Implement QueryCount() method on Client to get the number of DNS queries performed.

### Fixed
- Fix panic condition when list of domains to resolve is empty.

## [0.1.1] - 2021-04-05
### Fixed
- Fix race condition happening when creating new connections.

## [0.1.0] - 2021-03-18
### Added
- Initial implementation.

[Unreleased]: https://github.com/d3mondev/resolvermt/compare/v0.3.0...HEAD
[0.1.0]: https://github.com/d3mondev/resolvermt/releases/tag/v0.1.0
[0.1.1]: https://github.com/d3mondev/resolvermt/releases/tag/v0.1.1
[0.2.0]: https://github.com/d3mondev/resolvermt/releases/tag/v0.2.0
[0.3.0]: https://github.com/d3mondev/resolvermt/releases/tag/v0.3.0
