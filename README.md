# Terraform Provider: dotenv

![Terraform](https://img.shields.io/badge/Terraform-%3E%3D1.10-7B42BC?logo=terraform&logoColor=white)
![OpenTofu](https://img.shields.io/badge/OpenTofu-%3E%3D1.11-FFDA18?logo=opentofu&logoColor=black)
![Go](https://img.shields.io/badge/Go-%3E%3D1.26-00ADD8?logo=go&logoColor=white)
[![GitHub Release](https://img.shields.io/github/v/release/pfandie/terraform-provider-dotenv?logo=github&label=Release)](https://github.com/pfandie/terraform-provider-dotenv/releases/latest)
[![Tests](https://img.shields.io/github/actions/workflow/status/pfandie/terraform-provider-dotenv/test.yml?branch=main&logo=github&label=Tests)](https://github.com/pfandie/terraform-provider-dotenv/actions/workflows/test.yml)

Load environment variables from `.env` files into Terraform. \
This provider reads dotenv files and hands you their key/value pairs in two flavors:
* a regular data source
* an [ephemeral resource](https://developer.hashicorp.com/terraform/language/resources/ephemeral) (Terraform >= 1.10)

Tested with both **Terraform** and **OpenTofu**.

## Documentation

The full schema documentation lives in [docs](./docs) and on the
[Terraform Registry](https://registry.terraform.io/providers/pfandie/dotenv/latest/docs).

## Usage

### Data source

```hcl
terraform {
  required_providers {
    dotenv = {
      source  = "pfandie/dotenv"
      version = ">= 0.1.0"
    }
  }
}

data "dotenv_file" "config" {
  env_file       = ".env.production"
  exclude_envs   = ["LOCAL_ONLY_FLAG"]
  sensitive_envs = ["API_TOKEN"]
}

resource "aws_lambda_function" "this" {
  # ...
  environment {
    variables = merge(
      data.dotenv_file.config.values,
      data.dotenv_file.config.sensitive_values,
    )
  }
}
```

### Ephemeral resource

Values opened ephemerally are not persisted in the plan or state file —
ideal for secrets that are written to write-only arguments:

```hcl
ephemeral "dotenv_file" "secrets" {
  env_file = ".env.secrets"
}

resource "aws_secretsmanager_secret_version" "app" {
  secret_id                = aws_secretsmanager_secret.app.id
  secret_string_wo         = jsonencode(ephemeral.dotenv_file.secrets.values)
  secret_string_wo_version = 1
}
```

## Examples

- [Data source](./examples/data-sources/dotenv_file)
- [Ephemeral resource](./examples/ephemeral-resources/dotenv_file)

## Requirements

| Name                                                             | Version |
|------------------------------------------------------------------|---------|
| [Terraform](https://developer.hashicorp.com/terraform/downloads) | >= 1.10 |
| [OpenTofu](https://opentofu.org/docs/intro/install/)             | >= 1.11 |

## Resources

| Name                                              | Type               |
|---------------------------------------------------|--------------------|
| [dotenv_file](./docs/data-sources/file.md)        | data source        |
| [dotenv_file](./docs/ephemeral-resources/file.md) | ephemeral resource |

## Inputs

| Name           | Description                                                       | Type           | Default  | Required |
|----------------|-------------------------------------------------------------------|----------------|----------|:--------:|
| env_path       | Directory containing the env file                                 | `string`       | `"."`    |    no    |
| env_file       | Name of the env file inside `env_path`                            | `string`       | `".env"` |    no    |
| include_envs   | Only these keys will be included. Exclusive to `exclude_envs`     | `list(string)` | `null`   |    no    |
| exclude_envs   | These keys won't be included. Exclusive to `include_envs`         | `list(string)` | `null`   |    no    |
| sensitive_envs | These keys are returned in `sensitive_values` instead of `values` | `list(string)` | `null`   |    no    |
| include_empty  | Also load empty keys (e.g. `FOO=`)                                | `bool`         | `false`  |    no    |

## Outputs

| Name             | Description                                                       |
|------------------|-------------------------------------------------------------------|
| values           | All key/value pairs without `sensitive_envs`                      |
| sensitive_values | Key/value pairs defined in `sensitive_envs` (marked as sensitive) |

## Developing the Provider

### Requirements

- [Go](https://go.dev/doc/install) >= 1.26 (to build the provider plugin)

### Building and testing

```sh
# build the provider
make build

# unit tests
make test

# acceptance tests (runs real terraform plans against a local test binary)
make testacc

# acceptance tests against OpenTofu (requires a local tofu binary)
make testacc-tofu

# linting (requires golangci-lint v2)
make lint
```

The make targets are the editor-agnostic entry point — every IDE and the CI
use the same commands. On top of that there are thin wrappers for:

- **IntelliJ/GoLand**: shared run configurations in [.run](./.run) appear
  automatically:
  - `Unit Tests`
  - `Acceptance Tests (Terraform)`
  - `Acceptance Tests (OpenTofu)`
- **VS Code**: tasks in [.vscode/tasks.json](./.vscode/tasks.json) ("Run Task") wrap the same make targets, [launch.json](./.vscode/launch.json)
  provides debug configurations with breakpoints

### Trying the provider locally

Build and install the provider into your `GOBIN`:

```bash
make install
```

Then point Terraform at the local build via a
[`dev_overrides`](https://developer.hashicorp.com/terraform/cli/config/config-file#development-overrides-for-provider-developers)
block in `~/.terraformrc`:

```hcl
provider_installation {
  dev_overrides {
    "pfandie/dotenv" = "/Users/<you>/go/bin"
  }
  direct {}
}
```

With the override active, skip `terraform init` and run `terraform plan` /
`terraform apply` directly

### Generating docs

The documentation in [docs](./docs) is generated — do not edit it by hand:

```sh
make generate
```

### Releasing

Releases are automated with [release-please](https://github.com/googleapis/release-please):
commit messages follow [Conventional Commits](https://www.conventionalcommits.org/):
- `feat:`
- `fix:`
- `chore:`
- ...

## License

[Apache-2.0](./LICENSE) — Copyright 2026 Hans Mayer ([@pfandie](https://github.com/pfandie))
