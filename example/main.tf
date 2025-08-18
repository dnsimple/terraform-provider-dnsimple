terraform {
  required_version = ">= 1.11"

  required_providers {
    dnsimple = {
      source  = "dnsimple/dnsimple"
      version = ">= 1.9.0"
    }
  }
}

provider "dnsimple" {
  token   = var.dnsimple_token
  account = var.dnsimple_account
  sandbox = true
}


variable "dnsimple_token" {
  description = "DNSimple API Token"
  type        = string
  sensitive   = true
}

variable "dnsimple_account" {
  description = "DNSimple Account ID"
  type        = string
}


# Create a record.
resource "dnsimple_zone_record" "record_1755513796" {
  zone_name = "example.com"
  name      = "tf"
  value     = "Hello Terraform!"
  type      = "TXT"
}
