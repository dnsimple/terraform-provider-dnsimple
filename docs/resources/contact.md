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
  fax = "+1849491024"
  email = "apple@contoso.com"
}
```

## Argument Reference

The following argument(s) are supported:

* `label` - Label
* `first_name` - (Required) First name
* `last_name` - (Required) Last name
* `organization_name` - Organization name
* `job_title` - Job title
* `address1` - (Required) Address line 1
* `address2` - Address line 2
* `city` - (Required) City
* `state_province` - (Required) State province
* `postal_code` - (Required) Postal code
* `country` - (Required) Country
* `phone` - (Required) Phone
* `fax` - Fax
* `email` - (Required) Email

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
      "phone": "+18001234567",
      "fax": "+18011234567",
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
      "phone": "+18881234567",
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
