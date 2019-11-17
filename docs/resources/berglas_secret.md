# `berglas_secret` Resource

Create and manage secrets via [Berglas][berglas].

## Example Usage

Create a secret:

```hcl
resource "berglas_secret" "apikey" {
  bucket    = "my-bucket"
  name      = "my-apikey"
  key       = "projects/p/locations/l/keyRings/r/cryptoKeys/k"
  plaintext = "my-super-secret"
}
```

## Argument Reference

-   `bucket` - (Required) Name of the bucket where secrets are stored.

-   `name` - (Required) Name of the secret to create.

-   `key` - (Required) Fully-qualified [Cloud KMS][cloud-kms] key name to use
    for envelope encryption. This key must already exist.

-   `plaintext` - (Required) Secret material to store.

## Attribute Reference

-   `generation` - Version of the object.

-   `metageneration` - CAS token for concurrent operations (most users will
    never need this).

## Import

Import existing secrets by providing the id:

```sh
terraform import berglas_secret.api my-bucket/my-secret#generation
```

[berglas]: https://github.com/GoogleCloudPlatform/berglas
[cloud-kms]: https://cloud.google.com/kms
