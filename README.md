# Terraform Berglas Provider

This is a [Terraform][terraform] provider for interacting with
[Berglas][berglas].


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


## Usage

1. Create a Terraform configuration file:

    ```hcl
    resource "berglas_secret" "demo" {
      bucket    = "my-bucket"
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

For more examples, please see the [examples][examples] folder in this
repository.

## Reference

See [the documentation](doc/).

## FAQ

**Q: Is it secure?**<br>
A: Secrets will be stored in plaintext in the Terraform state. You should only
use this with remote state. Please see [sensitive state][sensitive-state] for
more information.


[berglas]: https://github.com/GoogleCloudPlatform/berglas
[terraform]: https://www.terraform.io/
[releases]: https://github.com/sethvargo/terraform-provider-berglas/releases
[sensitive-state]: https://www.terraform.io/docs/extend/best-practices/sensitive-state.html
