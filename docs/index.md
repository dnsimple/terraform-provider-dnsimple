---
page_title: "Provider: DNSimple"
---

# DNSimple Provider

The DNSimple provider allows you to manage DNS records, domains, certificates, and other DNSimple resources using Terraform.

This provider enables you to treat your DNS and domain infrastructure as code, making it easier to version, review, and manage your DNSimple resources alongside your other infrastructure.

[![IMAGE_ALT](https://img.youtube.com/vi/cTWP1MWA-0c/0.jpg)](https://www.youtube.com/watch?v=cTWP1MWA-0c)

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.12
- DNSimple account with API access

## Installation

Add the DNSimple provider to your Terraform configuration:

```hcl
terraform {
  required_version = ">= 1.12"

  required_providers {
    dnsimple = {
      source  = "dnsimple/dnsimple"
      version = "~> 1.9"
    }
  }
}
```

Then run `terraform init` to download the provider.

## Authentication

The provider requires authentication credentials to interact with the DNSimple API. You can provide credentials in several ways:

1. **Provider configuration** (recommended for development)
2. **Environment variables** (recommended for CI/CD and production)

### Using Provider Configuration

```hcl
provider "dnsimple" {
  token   = var.dnsimple_token
  account = var.dnsimple_account
}
```

### Using Environment Variables

```bash
export DNSIMPLE_TOKEN="your-api-token"
export DNSIMPLE_ACCOUNT="your-account-id"
```

See the [Argument Reference](#argument-reference) section below for all configuration options.

## Example Usage

Configure the provider:

```hcl
terraform {
  required_version = ">= 1.12"

  required_providers {
    dnsimple = {
      source  = "dnsimple/dnsimple"
      version = "~> 1.9"
    }
  }
}

provider "dnsimple" {
  token   = var.dnsimple_token
  account = var.dnsimple_account
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
```

Now use the available resources to perform actions like managing DNS records or registering domains.

To manage your DNS records:

```hcl
# Create a zone
resource "dnsimple_zone" "example" {
  name = "example.com"
}

# Create DNS records
resource "dnsimple_zone_record" "www" {
  zone_name = dnsimple_zone.example.name
  name      = "www"
  value     = "192.0.2.1"
  type      = "A"
  ttl       = 3600
}

resource "dnsimple_zone_record" "apex" {
  zone_name = dnsimple_zone.example.name
  name      = ""
  value     = "192.0.2.1"
  type      = "A"
  ttl       = 3600
}
```

To register a domain:

```hcl
# Create a contact for domain registration
resource "dnsimple_contact" "registrant" {
  label           = "Main Contact"
  first_name      = "John"
  last_name       = "Doe"
  organization_name = "Example Inc"
  address1        = "123 Main Street"
  city            = "San Francisco"
  state_province  = "California"
  postal_code     = "94105"
  country         = "US"
  phone           = "+1.4155551234"
  email           = "john@example.com"
}

# Register a domain
resource "dnsimple_registered_domain" "example_com" {
  name       = "example.com"
  contact_id = dnsimple_contact.registrant.id

  auto_renew_enabled    = true
  whois_privacy_enabled = true
  transfer_lock_enabled = true
}
```

For more elaborate use cases, and to learn more about the capabilities offered by the DNSimple Terraform provider, view the individual resource and data source pages.

## Argument Reference

The following arguments are supported in the provider configuration:

- **`token`** (Required) - The DNSimple [API v2 token](https://support.dnsimple.com/articles/api-access-token/). Can be provided via the `DNSIMPLE_TOKEN` environment variable. You can use either a User or Account token, but an Account token is recommended for better security and access control.

- **`account`** (Required) - The ID of the account associated with the token. Can be provided via the `DNSIMPLE_ACCOUNT` environment variable.

- **`sandbox`** (Optional) - Set to `true` to connect to the API [sandbox environment](https://developer.dnsimple.com/sandbox/) for testing. Can be provided via the `DNSIMPLE_SANDBOX` environment variable. Defaults to `false`.

- **`prefetch`** (Optional) - Set to `true` to enable prefetching zone records when dealing with large configurations. This is useful when you are dealing with API rate limitations given your number of zones and zone records. Can be provided via the `DNSIMPLE_PREFETCH` environment variable. Defaults to `false`.

- **`user_agent`** (Optional) - Custom string to append to the user agent used for sending HTTP requests to the API. Useful for identifying your automation or integration.

## Getting Help

- [Support article](https://support.dnsimple.com/articles/terraform-provider/) - Official support documentation
- [Developer API documentation](https://developer.dnsimple.com/) - Complete API reference
- [GitHub Repository](https://github.com/dnsimple/terraform-provider-dnsimple) - Source code and issue tracker

## Related Articles

- [Introducing DNSimple's Terraform Provider](https://blog.dnsimple.com/2021/12/introducing-dnsimple-terraform-provider/)
- [DNSimple, Terraform & Sentinel — A Guide to Policy as Code](https://blog.dnsimple.com/2023/05/policy-as-code/)
- [Manage Domain Transfer Locking and Contacts in Terraform](https://blog.dnsimple.com/2023/06/terraform-domain-registrations/)
- [How We Manage Domain and DNS Management with Infrastructure as Code](https://blog.dnsimple.com/2025/11/managing-domains-terraform-dnsimple/)
