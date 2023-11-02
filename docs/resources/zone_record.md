---
page_title: "DNSimple: dnsimple_zone_record"
---

# dnsimple\_zone\_record

Provides a DNSimple zone record resource.

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

* `zone_name` - (Required) The zone name to add the record to
* `name` - (Required) The name of the record
* `value` - (Required) The value of the record
* `type` - (Required) The type of the record
* `ttl` - (Optional) The TTL of the record - defaults to 3600
* `priority` - (Optional) The priority of the record - only useful for some record types
* `regions` - (Optional) A list of regions to serve the record from. You can find a list of supported values in our [developer documentation](https://developer.dnsimple.com/v2/zones/records/).


## Attributes Reference

* `id` - The record ID
* `zone_id` - The zone ID of the record
* `qualified_name` - The FQDN of the record
* `value_normalized` - The normalized value of the record

## Import

DNSimple resources can be imported using their parent zone name (domain name) and numeric record ID.

**Importing record example.com with record ID 1234**

```bash
terraform import dnsimple_zone_record.resource_name example.com_1234
```

**Importing record www.example.com with record ID 1234**

```bash
terraform import dnsimple_zone_record.resource_name example.com_1234
```

The record ID can be found in the URL when editing a record on the DNSimple web dashboard.
