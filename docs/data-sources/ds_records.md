---
page_title: "DNSimple: dnsimple_ds_records"
---

# dnsimple\_ds\_records

Get the delegation signer (DS) records of a DNSSEC-signed domain hosted at DNSimple.

Use this data source when the zone is hosted at DNSimple but the domain is registered elsewhere: it returns the DS records you need to enter at your registrar. During a key-signing-key (KSK) rollover the domain temporarily has two DS records, and this data source returns all of them so the configuration does not go stale.

## Example Usage

```hcl
data "dnsimple_ds_records" "example" {
  domain = "example.com"
}

# Copy each DS record to the registrar
resource "example_registrar_ds_record" "example" {
  for_each = { for record in data.dnsimple_ds_records.example.records : record.key_tag => record }

  domain      = "example.com"
  key_tag     = each.value.key_tag
  algorithm   = each.value.algorithm
  digest_type = each.value.digest_type
  digest      = each.value.digest
}
```

## Argument Reference

The following arguments are supported:

- `domain` - (Required) The name or ID of the domain.

## Attributes Reference

The following attributes are exported:

- `id` - The domain name or ID.
- `records` - (List) A list of delegation signer records for the domain. Empty when the domain is not DNSSEC-signed. (see [below for nested schema](#nested-schema-for-records))

### Nested Schema for `records`

Attributes Reference:

- `id` (Number) - The DS record ID in DNSimple.
- `algorithm` (String) - The DNSSEC algorithm number used to generate the signature, e.g. `13` for ECDSA P256 SHA256.
- `digest` (String) - The digest of the DS record.
- `digest_type` (String) - The digest type of the DS record, e.g. `2` for SHA-256.
- `key_tag` (String) - The key tag of the DS record.
- `public_key` (String) - The public key of the DS record.
- `created_at` (String) - When the DS record was created at DNSimple.
- `updated_at` (String) - When the DS record was last updated at DNSimple.
