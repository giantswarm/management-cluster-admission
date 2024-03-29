# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.8.1] - 2023-01-27

- Add ignores for Nancy.

## [0.8.0] - 2022-12-08

### Fixed

- Exclude `crossplane` which consists of many deployments.

## [0.7.0] - 2022-03-31

### Changed

- Remove `apiextensions` dependency.
- Use `v1beta1` CAPI CRDs.

## [0.6.3] - 2022-02-01

### Fixed

- In `DeploymentValidator`, when checking if the Deployment is unique, check `app.kubernetes.io/component` label if it is present.

## [0.6.2] - 2021-09-10

### Fixed

- Fix manifests after changes introduced in v0.6.1.

## [0.6.1] - 2021-09-09

### Fixed

- Exclude `flux-app` which consists of many deployments.

## [0.6.0] - 2021-08-11

### Changed

- Migrate to configuration management.
- Update `architect-orb` to v4.0.0.

## [0.5.0] - 2021-08-05

### Changed

- Fetch CAPI `v1alpha3` `Cluster` resources, instead of `v1alpha2` for validating if an organization can be deleted or not.

## [0.4.0] - 2021-05-13

### Added

- Allow colliding `app-operator` apps from outside `giantswarm` namespace.

## [0.3.0] - 2021-05-10

### Changed

- Fetch CAPI `v1alpha2` `Cluster` resources, instead of `v1alpha3` for validating if an organization can be deleted or not.

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

- Initial release

[Unreleased]: https://github.com/giantswarm/management-cluster-admission/compare/v0.8.1...HEAD
[0.8.1]: https://github.com/giantswarm/management-cluster-admission/compare/v0.8.0...v0.8.1
[0.8.0]: https://github.com/giantswarm/management-cluster-admission/compare/v0.7.0...v0.8.0
[0.7.0]: https://github.com/giantswarm/management-cluster-admission/compare/v0.6.3...v0.7.0
[0.6.3]: https://github.com/giantswarm/management-cluster-admission/compare/v0.6.2...v0.6.3
[0.6.2]: https://github.com/giantswarm/management-cluster-admission/compare/v0.6.1...v0.6.2
[0.6.1]: https://github.com/giantswarm/management-cluster-admission/compare/v0.6.0...v0.6.1
[0.6.0]: https://github.com/giantswarm/management-cluster-admission/compare/v0.5.0...v0.6.0
[0.5.0]: https://github.com/giantswarm/management-cluster-admission/compare/v0.4.0...v0.5.0
[0.4.0]: https://github.com/giantswarm/management-cluster-admission/compare/v0.3.0...v0.4.0
[0.3.0]: https://github.com/giantswarm/management-cluster-admission/compare/v0.2.1...v0.3.0
[0.2.1]: https://github.com/giantswarm/management-cluster-admission/compare/v0.2.0...v0.2.1
[0.2.0]: https://github.com/giantswarm/management-cluster-admission/compare/v0.1.2...v0.2.0
[0.1.2]: https://github.com/giantswarm/management-cluster-admission/compare/v0.1.1...v0.1.2
[0.1.1]: https://github.com/giantswarm/management-cluster-admission/compare/v0.1.0...v0.1.1
[0.1.0]: https://github.com/giantswarm/management-cluster-admission/releases/tag/v0.1.0
