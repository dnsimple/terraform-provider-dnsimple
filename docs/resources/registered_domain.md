---
page_title: "DNSimple: dnsimple_registered_domain"
---

# dnsimple\_registered\_domain

Provides a DNSimple registered domain resource.

## Example Usage

```hcl
resource "dnsimple_contact" "alice_main" {
  label = "Alice Appleseed"
  first_name = "Alice Main"
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

resource "dnsimple_registered_domain" "appleseed_bio" {
  name = "appleseed.bio"

  contact_id            = dnsimple_contact.alice_main.id
  auto_renew_enabled    = true
  transfer_lock_enabled = true
  whois_privacy_enabled = true
  dnssec_enabled        = false

  extended_attributes = {
    "bio_agree" = "I Agree"
  }
}
```

## Argument Reference

The following argument(s) are supported:

* `name` - (Required) The domain name to be registered
* `contact_id` - (Required) The ID of the contact to be used for the domain registration. The contact ID can be changed after the domain has been registered. The change will result in a new registrant change this may result in a [60-day lock](https://support.dnsimple.com/articles/icann-60-day-lock-registrant-change/).
* `auto_renew_enabled` - (Optional) Whether the domain should be set to auto-renew (default: `false`)
* `whois_privacy_enabled` - (Optional) Whether the domain should have WhoIs privacy enabled (default: `false`)
* `dnssec_enabled` - (Optional) Whether the domain should have DNSSEC enabled (default: `false`)
* `transfer_lock_enabled` - (Optional) Whether the domain transfer lock protection is enabled (default: `true`)
* `premium_price` - (Optional) The premium price for the domain registration. This is only required if the domain is a premium domain. You can use our [Check domain API](https://developer.dnsimple.com/v2/registrar/#checkDomain) to check if a domain is premium. And [Retrieve domain prices API](https://developer.dnsimple.com/v2/registrar/#getDomainPrices) to retrieve the premium price for a domain.
* `extended_attributes` - (Optional) A map of extended attributes to be set for the domain registration. To see if there are any required extended attributes for any TLD use our [Lists the TLD Extended Attributes API](https://developer.dnsimple.com/v2/tlds/#getTldExtendedAttributes). The values provided in the `extended_attributes` will also be sent when a registrant change is initiated as part of changing the `contact_id`.
* `timeouts` - (Optional) (see [below for nested schema](#nestedblock--timeouts))

# Attributes Reference

- `id` - The ID of this resource.
- `unicode_name` - The domain name in Unicode format.
- `state` - The state of the domain.
- `domain_registration` - (Block) The domain registration details. (see [below for nested schema](#nestedblock--domain_registration))

<a id="nestedblock--timeouts"></a>

### Nested Schema for `timeouts`

Optional:

- `create` (String) - The timeout for the read operation e.g. `5m`
- `update` (String) - The timeout for the read operation e.g. `5m`

<a id="nestedblock--domain_registration"></a>

### Nested Schema for `domain_registration`

Attributes Reference:

- `id` (Number) - The ID of the domain registration.
- `state` (String) - The state of the domain registration.
- `period` (Number) - The registration period in years.

## Import

DNSimple registered domains can be imported using their domain name and **optionally** with domain registration ID.

**Importing registered domain example.com**

```bash
terraform import dnsimple_registered_domain.resource_name example.com
```

**Importing registered domain example.com with domain registration ID 1234**

```bash
terraform import dnsimple_registered_domain.resource_name example.com_1234
```

~> **Note:** At present there is no way to retrieve the domain registration ID from the DNSimple API or UI. You will need to have noted the ID when you created the domain registration. Prefer using the domain name only when importing.
