---
page_title: "DNSimple: dnsimple_registered_domain_contact"
---

# dnsimple\_registered\_domain\_contact

Provides a way to manage your registered domain's contact at DNSimple.

-> **Note:** The registrant change API is currently in developer preview and is subject to change.

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

resource "dnsimple_registered_domain_contact" "appleseed_bio" {
  domain_id  = "appleseed.bio"
  contact_id = dnsimple_contact.alice_main.id

  extended_attributes = {
    "bio_agree" = "I Agree"
  }
}
```

## Argument Reference

The following argument(s) are supported:

* `domain_id` - (Required) The domain name or id for which the contact should be updated
* `contact_id` - (Required) The ID of the new contact
* `extended_attributes` - (Optional) A map of extended attributes to be set for the domain registration. To see if there are any required extended attributes for the contact change. You can use the `dnsimple_registrant_change_check` data source to check if there are any required extended attributes for the contact change. (see [data source](../data-sources/registrant_change_check.md))
* `timeouts` - (Optional) (see [below for nested schema](#nestedblock--timeouts))

# Attributes Reference

- `id` - The ID of this resource.
- `account_id` - The account ID that owns the domain.
- `state` - The state of the registrant change.
* `registry_owner_change` - (Boolean) Whether the registrant change has resulted in an owner change at the registry.
* `irt_lock_lifted_by` - (String) The date when the Inter-Registrar Transfer (IRT) lock will/was lifted.


<a id="nestedblock--timeouts"></a>

### Nested Schema for `timeouts`

Optional:

- `create` (String) - The timeout for the read operation e.g. `5m`
- `update` (String) - The timeout for the read operation e.g. `5m`

## Import

DNSimple registrant change resource can be imported using its ID.

**Importing registrant change with ID 1234**

```bash
terraform import dnsimple_registered_domain_contact.resource_name 1234
```
