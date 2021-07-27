terraform {
  required_providers {
    dnsimple = {
      source = "dnsimple/dnsimple"
      version = "0.5.3-10-g35e3384"
    }
  }
}

provider "dnsimple" {
    token = var.dnsimple_token
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


variable "dnsimple_domain" {
  description = "DNSimple Domain"
  type        = string
}

# Create a record
resource "dnsimple_zone_record" "test-txt" {
    zone_name = var.dnsimple_domain
    name   = "test-tf-txt"
    value  = "Hello Terraform!"
    type   = "TXT"
    ttl    = 3600
}