# `berglas_secret` Data Source

Access [Berglas][berglas] secrets.

## Example Usage

Access a secret:

```hcl
data "berglas_secret" "apikey" {
  bucket = "my-bucket"
  name   = "my-apikey"
}

// data.berglas_secret.apikey.plaintext contains the plaintext secret
```

## Argument Reference

-   `bucket` - (Required) Name of the bucket where secrets are stored.

-   `name` - (Required) Name of the secret to create.

-   `generation` (Optional) Specific version to access.

## Attribute Reference

-   `key` - Fully-qualified [Cloud KMS][cloud-kms] key name used for envelope
    encryption.

-   `plaintext` - Secret material.

-   `metageneration` - CAS token for concurrent operations (most users will
    never need this).

[berglas]: https://github.com/GoogleCloudPlatform/berglas
[cloud-kms]: https://cloud.google.com/kms
