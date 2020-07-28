---
layout: "dnsimple"
page_title: "DNSimple: dnsimple_email_forward"
sidebar_current: "docs-dnsimple-resource-email-forward"
description: |-
  Provides a DNSimple email forward resource.
---

# dnsimple\_email\_forward

Provides a DNSimple email forward resource.

## Example Usage

```hcl
# Add an email forwarding rule to the domain
resource "dnsimple_email_forward" "foobar" {
  domain = "${var.dnsimple_domain}"
  from   = "sales"
  to     = "jane.doe@example.com"
}
```

## Argument Reference

The following arguments are supported:

* `domain` - (Required) The domain to add the email forwarding rule to
* `from` - (Required) The source email address on the domain
* `to` - (Required) The destination email address on another domain

## Attributes Reference

The following attributes are exported:

* `id` - The email forward ID
* `from` - The source email address on the domain
* `to` - The destination email address on another domain
