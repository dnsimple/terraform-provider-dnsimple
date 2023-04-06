---
page_title: "Provider: DNSimple"
---

# DNSimple Provider

The DNSimple provider is used to interact with the resources supported by DNSimple. The provider needs to be configured
with the proper credentials before it can be used.

Use the navigation to the left to read about the available resources.

[![IMAGE_ALT](https://img.youtube.com/vi/cTWP1MWA-0c/0.jpg)](https://www.youtube.com/watch?v=cTWP1MWA-0c)

## Example Usage

```hcl
# Configure the DNSimple provider
provider "dnsimple" {
  token = "${var.dnsimple_token}"
  account = "${var.dnsimple_account}"
}

# Create a record
resource "dnsimple_zone_record" "www" {
  # ...
}

# Create an email forward
resource "dnsimple_email_forward" "hello" {
  # ...
}
```


## Argument Reference

The following arguments are supported:

* `token` - (Required) The DNSimple API v2 token. It must be provided, but it can also be sourced from the `DNSIMPLE_TOKEN` environment variable. Please note that this must be an [API v2 token](https://support.dnsimple.com/articles/api-access-token/). You can use either an User or Account token, but an Account token is recommended.
* `account` - (Required) The ID of the account associated with the token. It must be provided, but it can also be sourced from the `DNSIMPLE_ACCOUNT` environment variable.
* `sandbox` - Set to true to connect to the API [sandbox environment](https://developer.dnsimple.com/sandbox/). `DNSIMPLE_SANDBOX` environment variable can also be used.
* `prefetch` - Set to true to enable prefetching `ZoneRecords` when dealing with large configurations. This is useful
when you are dealing with API rate limitations given your number of zones and zone records. `DNSIMPLE_PREFETCH` environment variable can also be used.
* `user_agent` - (Optional) Custom string to append to the user agent used for sending HTTP requests to the API.

## Helpful Links

* [Blog post](https://blog.dnsimple.com/2021/12/introducing-dnsimple-terraform-provider/)
* [Support article](https://support.dnsimple.com/articles/terraform-provider/)
* [GitHub Repo](https://github.com/dnsimple/terraform-provider-dnsimple)
