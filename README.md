# Terraform Berglas Provider

This is a [Terraform][terraform] provider for interacting with
[Berglas][berglas].

**Secrets will be stored in plaintext in the Terraform state. You should only
use this with provider with Terraform remote state. For more information, please
see [sensitive state][sensitive-state].**


## Installation

1. Download the latest compiled binary from [GitHub releases][releases].

1. Unzip/untar the archive.

1. Move it into `$HOME/.terraform.d/plugins`:

    ```sh
    $ mkdir -p $HOME/.terraform.d/plugins
    $ mv terraform-provider-berglas $HOME/.terraform.d/plugins/terraform-provider-berglas
    ```

1. Create your Terraform configurations as normal, and run `terraform init`:

    ```sh
    $ terraform init
    ```

    This will find the plugin locally.

1. If you haven't already, [bootstrap berglas](https://github.com/GoogleCloudPlatform/berglas#setup)

#### Optionally

If using terraform v0.13+ you can create a `versions.tf` file to pull the plugin during `terraform init` without installing it locally:

```hcl
terraform {
  required_providers {
    berglas = {
      source  = "sethvargo/berglas"
      version = "~> 0.1"
    }
  }
}
```


## Usage

1. Create a Terraform configuration file:

    ```hcl
    resource "berglas_secret" "demo" {
      bucket    = "my-bucket"
      key       = "projects/${var.project_id}/locations/global/keyRings/berglas/cryptoKeys/berglas-key"
      name      = "demo"
      plaintext = "p@s$w0rd!"
    }
    ```

1. Run `terraform init` to pull in the provider:

    ```sh
    $ terraform init
    ```

1. Run `terraform plan` and `terraform apply`:

    ```sh
    $ terraform plan

    $ terraform apply
    ```

## Examples

For more examples, please see the [examples](examples/) folder in this
repository.

## Reference

See [the documentation](docs/).


[berglas]: https://github.com/GoogleCloudPlatform/berglas
[terraform]: https://www.terraform.io/
[releases]: https://github.com/sethvargo/terraform-provider-berglas/releases
[sensitive-state]: https://www.terraform.io/docs/extend/best-practices/sensitive-state.html
