---
page_title: "DNSimple: dnsimple_contact"
---

# dnsimple\_contact

Provides a DNSimple contact resource.

## Example Usage

```hcl
# Create a contact
resource "dnsimple_contact" "me" {
  label = "Apple Appleseed"
  first_name = "Apple"
  last_name = "Appleseed"
  organization_name = "Contoso"
  job_title = "Manager"
  address1 = "Level 1, 2 Main St"
  address2 = "Marsfield"
  city = "San Francisco"
  state_province = "California"
  postal_code = "90210"
  country = "US"
  phone = "+1401239523"
  fax = "+1.849491024"
  email = "apple@contoso.com"
}
```

## Argument Reference

The following argument(s) are supported:

* `label` - (String) A descriptive label for the contact to help identify it.
* `first_name` - (Required, String) The first name of the contact person.
* `last_name` - (Required, String) The last name of the contact person.
* `organization_name` - (String) The name of the organization or company associated with the contact.
* `job_title` - (String) The job title or position of the contact person within the organization.
* `address1` - (Required, String) The primary address line (street address, building number, etc.).
* `address2` - (String) The secondary address line (apartment, suite, floor, etc.).
* `city` - (Required, String) The city where the contact is located.
* `state_province` - (Required, String) The state, province, or region where the contact is located.
* `postal_code` - (Required, String) The postal code, ZIP code, or equivalent for the contact's location.
* `country` - (Required, String) The two-letter ISO country code (e.g., "US", "CA", "IT") for the contact's location.
* `phone` - (Required, String) The contact's phone number. Use international format with country code (e.g., "+1.4012345678" for US numbers).
* `fax` - (String) The contact's fax number. Use international format with country code (e.g., "+1.8491234567" for US numbers).
* `email` - (Required, String) The contact's email address.

# Attributes Reference

- `id` - The ID of this resource.
- `account_id` - The account ID for the contact.
- `phone_normalized` - The phone number, normalized.
- `fax_normalized` - The fax number, normalized.
- `created_at` - Timestamp representing when this contact was created.
- `updated_at` - Timestamp representing when this contact was updated.


## Import

DNSimple contacts can be imported using their numeric ID.

```bash
terraform import dnsimple_contact.resource_name 5678
```

The ID can be found within [DNSimple Contacts API](https://developer.dnsimple.com/v2/contacts/#listContacts). Check out [Authentication](https://developer.dnsimple.com/v2/#authentication) in API Overview for available options.

```bash
curl -u 'EMAIL:PASSWORD' https://api.dnsimple.com/v2/1234/contacts?label_like=example.com | jq
{
  "data": [
    {
      "id": 1,
      "account_id": 1010,
      "label": "Default",
      "first_name": "First",
      "last_name": "User",
      "job_title": "CEO",
      "organization_name": "Awesome Company",
      "email": "first@example.com",
      "phone": "+1.8001234567",
      "fax": "+1.8011234567",
      "address1": "Italian Street, 10",
      "address2": "",
      "city": "Roma",
      "state_province": "RM",
      "postal_code": "00100",
      "country": "IT",
      "created_at": "2013-11-08T17:23:15Z",
      "updated_at": "2015-01-08T21:30:50Z"
    },
    {
      "id": 2,
      "account_id": 1010,
      "label": "",
      "first_name": "Second",
      "last_name": "User",
      "job_title": "",
      "organization_name": "",
      "email": "second@example.com",
      "phone": "+1.8881234567",
      "fax": "",
      "address1": "French Street",
      "address2": "c/o Someone",
      "city": "Paris",
      "state_province": "XY",
      "postal_code": "00200",
      "country": "FR",
      "created_at": "2014-12-06T15:46:18Z",
      "updated_at": "2014-12-06T15:46:18Z"
    }
  ],
  "pagination": {
    "current_page": 1,
    "per_page": 30,
    "total_entries": 2,
    "total_pages": 1
  }
}
```
