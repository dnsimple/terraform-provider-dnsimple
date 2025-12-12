---
page_title: "DNSimple: dnsimple_contact"
---

# dnsimple\_contact

Provides a DNSimple contact resource.

## Example Usage

```hcl
resource "dnsimple_contact" "example" {
  label            = "Main Contact"
  first_name       = "John"
  last_name        = "Doe"
  organization_name = "Example Inc"
  job_title        = "Manager"
  address1         = "123 Main Street"
  address2         = "Suite 100"
  city             = "San Francisco"
  state_province   = "California"
  postal_code      = "94105"
  country          = "US"
  phone            = "+1.4155551234"
  fax              = "+1.4155555678"
  email            = "john@example.com"
}
```

## Argument Reference

The following arguments are supported:

- `label` - (Optional) A descriptive label for the contact to help identify it.
- `first_name` - (Required) The first name of the contact person.
- `last_name` - (Required) The last name of the contact person.
- `organization_name` - (Optional) The name of the organization or company associated with the contact.
- `job_title` - (Optional) The job title or position of the contact person within the organization.
- `address1` - (Required) The primary address line (street address, building number, etc.).
- `address2` - (Optional) The secondary address line (apartment, suite, floor, etc.).
- `city` - (Required) The city where the contact is located.
- `state_province` - (Required) The state, province, or region where the contact is located.
- `postal_code` - (Required) The postal code, ZIP code, or equivalent for the contact's location.
- `country` - (Required) The two-letter ISO country code (e.g., "US", "CA", "IT") for the contact's location.
- `phone` - (Required) The contact's phone number. Use international format with country code (e.g., "+1.4012345678" for US numbers).
- `fax` - (Optional) The contact's fax number. Use international format with country code (e.g., "+1.8491234567" for US numbers).
- `email` - (Required) The contact's email address.

## Attributes Reference

- `id` - The ID of this resource.
- `account_id` - The account ID for the contact.
- `phone_normalized` - The phone number, normalized.
- `fax_normalized` - The fax number, normalized.
- `created_at` - Timestamp representing when this contact was created.
- `updated_at` - Timestamp representing when this contact was updated.


## Import

DNSimple contacts can be imported using their numeric ID.

```bash
terraform import dnsimple_contact.example 5678
```

The contact ID can be found within the [DNSimple Contacts API](https://developer.dnsimple.com/v2/contacts/#listContacts). Check out [Authentication](https://developer.dnsimple.com/v2/#authentication) in API Overview for available options.
