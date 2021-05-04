# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.2.1] - 2021-05-04

### Fixed

- Increase memory limits to avoid OOMs.

## [0.2.0] - 2021-04-20

### Added

- Add webhook to validate deletion of organization CRs.

### Fixed

- Push also to vmware app collection.

## [0.1.2] - 2021-04-12

### Fixed

- Use immutable labels for selectors.
- Route alerts to Team Biscuit.

## [0.1.1] - 2021-04-12

### Fixed

- Exclude `vertical-pod-autoscaler-app` which consists of many deployments.

## [0.1.0] - 2021-04-09

[Unreleased]: https://github.com/giantswarm/management-cluster-admission/compare/v0.2.1...HEAD
[0.2.1]: https://github.com/giantswarm/management-cluster-admission/compare/v0.2.0...v0.2.1
[0.2.0]: https://github.com/giantswarm/management-cluster-admission/compare/v0.1.2...v0.2.0
[0.1.2]: https://github.com/giantswarm/management-cluster-admission/compare/v0.1.1...v0.1.2
[0.1.1]: https://github.com/giantswarm/management-cluster-admission/compare/v0.1.0...v0.1.1
[0.1.0]: https://github.com/giantswarm/management-cluster-admission/releases/tag/v0.1.0
