---
page_title: "DNSimple: dnsimple_domain"
---

# dnsimple\_domain

Provides a DNSimple domain resource.

## Example Usage

```hcl
# Create a domain
resource "dnsimple_domain" "foobar" {
  name = "${var.dnsimple.domain}"
}
```

## Argument Reference

The following argument(s) are supported:

* `name` - (Required) The domain name to be created