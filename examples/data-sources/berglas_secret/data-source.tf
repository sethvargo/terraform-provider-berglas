variable "bucket" {
  type = string
}

data "berglas_secret" "apikey" {
  bucket = var.bucket
  name   = "service-apikey"
}

output "demo" {
  value = data.berglas_secret.apikey.plaintext
}
