---
page_title: "DNSimple: dnsimple_email_forward"
---

# dnsimple\_email\_forward

Provides a DNSimple email forward resource.

## Example Usage

```hcl
# Add an email forwarding rule to the domain
resource "dnsimple_email_forward" "foobar" {
  domain = "${var.dnsimple_domain}"
  alias_name        = "sales"
  destination_email = "jane.doe@example.com"
}
```

## Argument Reference

The following arguments are supported:

* `domain` - (Required) The domain to add the email forwarding rule to
* `alias_name` - The name part (the part before the @) of the source email address on the domain
* `destination_email` - (Required) The destination email address on another domain

## Attributes Reference

The following attributes are exported:

* `id` - The email forward ID
* `alias_name` - The name part (the part before the @) of the source email address on the domain
* `alias_email` - The source email address on the domain
* `destination_email` - The destination email address on another domain
