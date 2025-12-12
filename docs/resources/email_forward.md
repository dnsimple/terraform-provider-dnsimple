---
page_title: "DNSimple: dnsimple_email_forward"
---

# dnsimple\_email\_forward

Provides a DNSimple email forward resource.

## Example Usage

```hcl
resource "dnsimple_email_forward" "example" {
  domain            = "example.com"
  alias_name        = "sales"
  destination_email = "alice@example.com"
}
```

## Argument Reference

The following arguments are supported:

- `domain` - (Required) The domain name to add the email forwarding rule to.
- `alias_name` - (Required) The name part (the part before the @) of the source email address on the domain.
- `destination_email` - (Required) The destination email address.

## Attributes Reference

The following attributes are exported:

- `id` - The email forward ID.
- `alias_email` - The source email address on the domain, in full form. This is a computed attribute.

## Import

DNSimple email forwards can be imported using the domain name and numeric email forward ID in the format `domain_name_email_forward_id`.

```bash
terraform import dnsimple_email_forward.example example.com_1234
```

The email forward ID can be found via the [DNSimple Email Forwards API](https://developer.dnsimple.com/v2/email-forwards/#listEmailForwards).
