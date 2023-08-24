---
page_title: "DNSimple: dnsimple_registrant_change_check"
---

# dnsimple\_registrant_change_check

Get information on the requirements of a registrant change.

-> **Note:** The registrant change API is currently in developer preview and is subject to change.

# Example Usage

Get registrant change requirements for the `dnsimple.com` domain and the contact with ID `1234`:

```hcl
data "dnsimple_registrant_change_check" "example" {
    domain_id = "dnsimple.com"
    contact_id = "1234"
}
```

# Argument Reference

The following arguments are supported:

* `domain_id` - (Required) The name or ID of the domain.
* `contact_id` - (Required) The ID of the contact you are planning to change to.

# Attributes Reference

The following additional attributes are exported:

* `contact_id` - The ID of the contact you are planning to change to.
* `domain_id` - The name or ID of the domain.
* `extended_attributes` - (List) A list of extended attributes that are required for the registrant change. (see [below for nested schema](#nestedblock--extended_attributes))
* `registry_owner_change` - (Boolean) Whether the registrant change is going to result in an owner change at the registry.

<a id="nestedblock--extended_attributes"></a>

### Nested Schema for `extended_attributes`

Attributes Reference:

- `name` (String) - The name of the extended attribute. e.g. `x-au-registrant-id-type`
- `description` (String) - The description of the extended attribute.
- `required` (Boolean) - Whether the extended attribute is required.
- `options` (List) - A list of options for the extended attribute. (see [below for nested schema](#nestedblock--options))

<a id="nestedblock--options"></a>

### Nested Schema for `extended_attributes.options`

Attributes Reference:

- `title` (String) - The human readable title of the option. e.g. `Australian Company Number (ACN)`
- `value` (String) - The value of the option.
- `description` (String) - The description of the option.


