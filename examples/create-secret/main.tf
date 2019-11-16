variable "bucket" {
  type = "string"
}

variable "kms_key" {
  type = "string"
}

resource "berglas_secret" "apikey" {
  bucket    = "${var.bucket}"
  name      = "service-apikey"
  key       = "${var.kms_key}"
  plaintext = "${other_resource.thing}" // example
}
