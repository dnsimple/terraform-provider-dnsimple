---
page_title: "DNSimple: dnsimple_lets_encrypt_certificate"
---

# dnsimple\_lets_encrypt_certificate

Provides a DNSimple Let's Encrypt certificate resource.

## Example Usage

```hcl
resource "dnsimple_lets_encrypt_certificate" "example" {
  domain_id       = "example.com"
  name            = "www"
  auto_renew      = true
  alternate_names = ["docs.example.com", "status.example.com"]
}
```

## Argument Reference

The following arguments are supported:

- `domain_id` - (Required) The domain name or ID to issue the certificate for.
- `name` - (Required) The certificate name; use `""` for the root domain. Wildcard names are supported.
- `alternate_names` - (Optional) List of alternate names (SANs) for the certificate.
- `auto_renew` - (Required) Whether the certificate should auto-renew.
- `signature_algorithm` - (Optional) The signature algorithm to use for the certificate.
- `timeouts` - (Block, Optional) (see [below for nested schema](#nested-schema-for-timeouts))

## Attributes Reference

The following attributes are exported:

- `id` - The certificate ID.
- `years` - The number of years the certificate will last.
- `state` - The state of the certificate.
- `authority_identifier` - The identifying certification authority (CA).
- `csr` - The certificate signing request.
- `expires_at` - The datetime when the certificate will expire.
- `created_at` - The datetime when the certificate was created.
- `updated_at` - The datetime when the certificate was last updated.

### Nested Schema for `timeouts`

Optional:

- `read` (String) - The timeout for the read operation, e.g., `5m`.

## Import

DNSimple Let's Encrypt certificates can be imported using the domain name and certificate ID in the format `domain_name_certificate_id`.

```bash
terraform import dnsimple_lets_encrypt_certificate.example example.com_1234
```

The certificate ID can be found via the [DNSimple Certificates API](https://developer.dnsimple.com/v2/certificates/#listCertificates).
