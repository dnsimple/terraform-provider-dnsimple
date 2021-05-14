terraform {
  required_providers {
    dnsimple = {
      source = "dnsimple/dnsimple"
      version = "0.5.1"
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
resource "dnsimple_record" "test-txt" {
    domain = var.dnsimple_domain
    name   = "test-tf-txt"
    value  = "Hello Terraform!"
    type   = "TXT"
    ttl    = 3600
}