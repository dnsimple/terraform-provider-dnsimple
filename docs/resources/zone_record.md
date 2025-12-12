---
page_title: "DNSimple: dnsimple_zone_record"
---

# dnsimple\_zone\_record

Provides a DNSimple zone record resource.

## Example Usage

```hcl
# Add a record to the root domain
resource "dnsimple_zone_record" "apex" {
  zone_name = "example.com"
  name      = ""
  value     = "192.0.2.1"
  type      = "A"
  ttl       = 3600
}

# Add a record to a subdomain
resource "dnsimple_zone_record" "www" {
  zone_name = "example.com"
  name      = "www"
  value     = "192.0.2.1"
  type      = "A"
  ttl       = 3600
}

# Add an MX record
resource "dnsimple_zone_record" "mx" {
  zone_name = "example.com"
  name      = ""
  value     = "mail.example.com"
  type      = "MX"
  priority  = 10
  ttl       = 3600
}
```

## Argument Reference

The following arguments are supported:

- `zone_name` - (Required) The zone name to add the record to.
- `name` - (Required) The name of the record. Use `""` for the root domain.
- `value` - (Required) The value of the record.
- `type` - (Required) The type of the record (e.g., `A`, `AAAA`, `CNAME`, `MX`, `TXT`). **The record type must be specified in UPPERCASE.**
- `ttl` - (Optional) The TTL of the record. Defaults to `3600`.
- `priority` - (Optional) The priority of the record. Only used for certain record types (e.g., `MX`, `SRV`).
- `regions` - (Optional) A list of regions to serve the record from. You can find a list of supported values in our [developer documentation](https://developer.dnsimple.com/v2/zones/records/).


## Attributes Reference

- `id` - The record ID.
- `zone_id` - The zone ID of the record.
- `qualified_name` - The fully qualified domain name (FQDN) of the record.
- `value_normalized` - The normalized value of the record.

## Import

DNSimple zone records can be imported using the zone name and numeric record ID in the format `zone_name_record_id`.

**Importing record for example.com with record ID 1234:**

```bash
terraform import dnsimple_zone_record.example example.com_1234
```

The record ID can be found in the URL when editing a record on the DNSimple web dashboard, or via the [DNSimple Zone Records API](https://developer.dnsimple.com/v2/zones/records/#listZoneRecords).
