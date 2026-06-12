# Changelog

## [1.1.1](https://github.com/pfandie/terraform-provider-dotenv/compare/v1.1.0...v1.1.1) (2026-06-12)


### Bug Fixes

* **docs:** adds process env to description/docs, adds 'make latest: true' to goreleaser ([bcb0712](https://github.com/pfandie/terraform-provider-dotenv/commit/bcb0712e43cc8d079364946f9f4c0b2b5569cd54))

## [1.1.0](https://github.com/pfandie/terraform-provider-dotenv/compare/v1.0.0...v1.1.0) (2026-06-11)


### Features

* enable possibility to read native env-vars and env-vars from ci ([fcb10d6](https://github.com/pfandie/terraform-provider-dotenv/commit/fcb10d68497ec4aec7a74b2034b7959f9f2e0c14))
* enable possibility to read native env-vars and env-vars from ci ([a503974](https://github.com/pfandie/terraform-provider-dotenv/commit/a503974dfb8344efed59b41eb9d422580502dcbc))


### Bug Fixes

* make use of github token to let release-please act as user and s… ([75c3e67](https://github.com/pfandie/terraform-provider-dotenv/commit/75c3e6733b3c3570187061b7c4557b914eb6e860))
* make use of github token to let release-please act as user and start workflows ([2fd1697](https://github.com/pfandie/terraform-provider-dotenv/commit/2fd1697e288025e14a53f4e0def19f879450253c))

## 1.0.0 (2026-06-11)


### Features

* add generated docs ([bb5cf9c](https://github.com/pfandie/terraform-provider-dotenv/commit/bb5cf9cd06e421ab0648fb788748aa16bbd57c9f))
* add gh-workflow for tests, add setting for golangci and goreleaser ([3262a6f](https://github.com/pfandie/terraform-provider-dotenv/commit/3262a6f88e0e22df47502c9a162baa16ec25a07b))
* add test setup for vscode and intellij/goland ([d2287a6](https://github.com/pfandie/terraform-provider-dotenv/commit/d2287a6b840fbd0d0772f7a9900b1e76fae8285a))
* add tests with test files ([b9f6e88](https://github.com/pfandie/terraform-provider-dotenv/commit/b9f6e885b2f79a319b7f67d604d89eea0a4e7a74))
* adds example files ([cb58703](https://github.com/pfandie/terraform-provider-dotenv/commit/cb5870395d1f24a570c0f4a4442397ce88eee6be))
* adds parser, data source and ephemeral sources ([a8737c8](https://github.com/pfandie/terraform-provider-dotenv/commit/a8737c86008e3de0ce0609fb5b4aee23f41892be))
* initial commit ([a9f5485](https://github.com/pfandie/terraform-provider-dotenv/commit/a9f54853d5b2bc8f3b2ce5ef3bf79a5c59942f51))
* updates gitignore file, remove opentofu 1.10.x as ephemeral is only available from opentofu 1.11.x ([e497685](https://github.com/pfandie/terraform-provider-dotenv/commit/e497685743c6ea4ff859f7eab0d4c9c34bec898d))


### Bug Fixes

* adjust tests for expect "non empty plan" for opentofu ([facdbed](https://github.com/pfandie/terraform-provider-dotenv/commit/facdbedca362a9cb18340be370426b0b30b69fba))
* remove changelog ignore in PR ([916a88d](https://github.com/pfandie/terraform-provider-dotenv/commit/916a88d67db5881337944900988f0769ab3071c2))
* revert 'ExpectNonEmptyPlan' and raise minimum version of opentofu ([fc781a6](https://github.com/pfandie/terraform-provider-dotenv/commit/fc781a6d1c07fef2b4ccc419596617a07ad5e748))
* trigger initial registry release ([d86aae5](https://github.com/pfandie/terraform-provider-dotenv/commit/d86aae5de08cd6f6f1c996cd4fd3334ce6f33165))
* trigger initial registry release ([376aef4](https://github.com/pfandie/terraform-provider-dotenv/commit/376aef4dfb773e0ad60e3857c606a6d77a977dda))
* trigger initial registry release ([a8a752e](https://github.com/pfandie/terraform-provider-dotenv/commit/a8a752eddc52dfa059376cf3975120d0461a40e5))
