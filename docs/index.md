# Berglas Provider

This is a Terraform provider for interacting with [Berglas][berglas].

!> Secrets will be stored in plaintext in the Terraform state. You should only
use this with remote state. Please see [sensitive state][sensitive-state] for
more information.

## Example Usage

Creating secrets:

```hcl
resource "berglas_secret" "apikey" {
  bucket    = "my-bucket"
  name      = "my-apikey"
  key       = "projects/p/locations/l/keyRings/r/cryptoKeys/k"
  plaintext = "my-super-secret"
}
```

Accessing secrets:

```hcl
data "berglas_secret" "apikey" {
  bucket = "my-bucket"
  name   = "my-apikey"
}

// data.berglas_secret.apikey.plaintext contains the secret contents
```

## Argument Reference

None

[berglas]: https://github.com/GoogleCloudPlatform/berglas
[sensitive-state]: https://www.terraform.io/docs/extend/best-practices/sensitive-state.html
