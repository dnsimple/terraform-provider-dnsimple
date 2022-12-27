---
page_title: "DNSimple: dnsimple_zone_record"
---

# dnsimple\_zone\_record

Provides a DNSimple zone record resource.

## Deprecation warning

You can still use the _deprecated_ `dnsimple_record` configuration, but be aware that it will be removed in the
upcoming 1.0.0 release. To migrate existing resources, please refer to the [migration guide](https://registry.terraform.io/providers/dnsimple/dnsimple/latest/docs/guides/resource-migration).

## Example Usage

```hcl
# Add a record to the root domain
resource "dnsimple_zone_record" "foobar" {
  zone_name = "${var.dnsimple_domain}"
  name   = ""
  value  = "192.168.0.11"
  type   = "A"
  ttl    = 3600
}
```

```hcl
# Add a record to a sub-domain
resource "dnsimple_zone_record" "foobar" {
  zone_name = "${var.dnsimple_domain}"
  name   = "terraform"
  value  = "192.168.0.11"
  type   = "A"
  ttl    = 3600
}
```

## Argument Reference

The following arguments are supported:

* `zone_name` - (Required) The domain to add the record to
* `name` - (Required) The name of the record
* `value` - (Required) The value of the record
* `type` - (Required) The type of the record
* `ttl` - (Optional) The TTL of the record
* `priority` - (Optional) The priority of the record - only useful for some record types


## Attributes Reference

The following attributes are exported:

* `id` - The record ID
* `name` - The name of the record
* `value` - The value of the record
* `type` - The type of the record
* `ttl` - The TTL of the record
* `priority` - The priority of the record
* `zone_id` - The domain ID of the record
* `qualified_name` - The FQDN of the record

## Import

DNSimple resources can be imported using their parent zone name (domain name) and numeric record ID.

**Importing record example.com with record ID 1234**

```
$ terraform import dnsimple_zone_record.resource_name example.com_1234
```

**Importing record www.example.com with record ID 1234**

```
$ terraform import dnsimple_zone_record.resource_name example.com_1234
```

The record ID can be found in the URL when editing a record on the DNSimple web dashboard.
