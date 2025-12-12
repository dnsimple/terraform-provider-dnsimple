---
page_title: "DNSimple: dnsimple_ds_record"
---

# dnsimple\_ds\_record

Provides a DNSimple domain delegation signer record resource.

## Example Usage

```hcl
resource "dnsimple_ds_record" "example" {
  domain      = "example.com"
  algorithm   = "8"
  digest      = "6CEEA0117A02480216EBF745A7B690F938860074E4AD11AF2AC573007205682B"
  digest_type = "2"
  key_tag     = "12345"
}
```

## Argument Reference

The following arguments are supported:

- `domain` - (Required) The domain name or numeric ID to create the delegation signer record for.
- `algorithm` - (Required) DNSSEC algorithm number as a string.
- `digest` - (Optional) The hexadecimal representation of the digest of the corresponding DNSKEY record.
- `digest_type` - (Optional) DNSSEC digest type number as a string.
- `key_tag` - (Optional) A key tag that references the corresponding DNSKEY record.
- `public_key` - (Optional) A public key that references the corresponding DNSKEY record.

## Attributes Reference

- `id` - The ID of this resource.
- `created_at` - The timestamp when the DS record was created.
- `updated_at` - The timestamp when the DS record was last updated.

## Import

DNSimple DS records can be imported using the domain name and numeric record ID in the format `domain_name_record_id`.

```bash
terraform import dnsimple_ds_record.example example.com_5678
```

The record ID can be found within the [DNSimple DNSSEC API](https://developer.dnsimple.com/v2/domains/dnssec/#listDomainDelegationSignerRecords). Check out [Authentication](https://developer.dnsimple.com/v2/#authentication) in API Overview for available options.
