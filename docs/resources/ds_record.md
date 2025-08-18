---
page_title: "DNSimple: dnsimple_ds_record"
---

# dnsimple\_ds\_record

Provides a DNSimple domain delegation signer record resource.

## Example Usage

```hcl
resource "dnsimple_ds_record" "foobar" {
  domain = "${var.dnsimple.domain}"
  algorithm = "8"
  digest = "6CEEA0117A02480216EBF745A7B690F938860074E4AD11AF2AC573007205682B"
  digest_type = "2"
  key_tag = "12345"
}
```

## Argument Reference

The following argument(s) are supported:

- `domain` - (Required) The domain name or numeric ID to create the delegation signer record for.
- `algorithm` - (Required) DNSSEC algorithm number as a string.
- `digest` - (Optional) The hexidecimal representation of the digest of the corresponding DNSKEY record.
- `digest_type` - (Optional) DNSSEC digest type number as a string.
- `keytag` - (Optional) A keytag that references the corresponding DNSKEY record.
- `public_key` - (Optional) A public key that references the corresponding DNSKEY record.

## Attributes Reference

- `id` - The ID of this resource.
- `created_at` - The time the DS record was created at.
- `updated_at` - The time the DS record was last updated at.

## Import

DNSimple DS record resources can be imported using their domain ID and numeric record ID.

```bash
terraform import dnsimple_domain_ds_signer.resource_name example.com_5678
```

The record ID can be found within [DNSimple DNSSEC API](https://developer.dnsimple.com/v2/domains/dnssec/#listDomainDelegationSignerRecords). Check out [Authentication](https://developer.dnsimple.com/v2/#authentication) in API Overview for available options.

```bash
curl -u 'EMAIL:PASSWORD' https://api.dnsimple.com/v2/1010/domains/example.com/ds_records | jq
{
  "data": [
    {
      "id": 24,
      "domain_id": 1010,
      "algorithm": "8",
      "digest": "C1F6E04A5A61FBF65BF9DC8294C363CF11C89E802D926BDAB79C55D27BEFA94F",
      "digest_type": "2",
      "keytag": "44620",
      "public_key": null,
      "created_at": "2017-03-03T13:49:58Z",
      "updated_at": "2017-03-03T13:49:58Z"
    }
  ],
  "pagination": {
    "current_page": 1,
    "per_page": 30,
    "total_entries": 1,
    "total_pages": 1
  }
}
```
