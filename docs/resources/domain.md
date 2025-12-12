---
page_title: "DNSimple: dnsimple_domain"
---

# dnsimple\_domain

Provides a DNSimple domain resource.

## Example Usage

```hcl
resource "dnsimple_domain" "example" {
  name = "example.com"
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Required) The domain name to be created.

## Attributes Reference

- `id` - The ID of this resource.
- `account_id` - The account ID for the domain.
- `auto_renew` - Whether the domain is set to auto-renew.
- `private_whois` - Whether the domain has WhoIs privacy enabled.
- `registrant_id` - The ID of the registrant (contact) for the domain.
- `state` - The state of the domain.
- `unicode_name` - The domain name in Unicode format.

## Import

DNSimple domains can be imported using the domain name.

```bash
terraform import dnsimple_domain.example example.com
```

The domain name can be found within the [DNSimple Domains API](https://developer.dnsimple.com/v2/domains/#listDomains). Check out [Authentication](https://developer.dnsimple.com/v2/#authentication) in API Overview for available options.
