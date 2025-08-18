---
page_title: "DNSimple: dnsimple_email_forward"
---

# dnsimple\_email\_forward

Provides a DNSimple email forward resource.

## Example Usage

```hcl
resource "dnsimple_email_forward" "foobar" {
  domain            = "${var.dnsimple_domain.name}"
  alias_name        = "sales"
  destination_email = "alice.appleseed@example.com"
}
```

## Argument Reference

The following arguments are supported:

- `domain` - (Required) The domain name to add the email forwarding rule to
- `alias_name` - The name part (the part before the @) of the source email address on the domain
- `destination_email` - (Required) The destination email address

## Attributes Reference

The following additional attributes are exported:

- `id` - The email forward ID
- `alias_email` - The source email address on the domain, in full form. This is a computed attribute.

## Import

DNSimple resources can be imported using the domain name and numeric email forward ID.

**Importing email forward for example.com with email forward ID 1234**

```bash
terraform import dnsimple_email_forward.resource_name example.com_1234
```
