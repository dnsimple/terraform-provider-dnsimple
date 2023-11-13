---
page_title: "DNSimple: dnsimple_lets_encrypt_certificate"
---

# dnsimple\_lets_encrypt_certificate

Provides a DNSimple Let's Encrypt certificate resource.

## Example Usage

```hcl
resource "dnsimple_lets_encrypt_certificate" "foobar" {
	domain_id = "${var.dnsimple.domain_id}"
	auto_renew = false
	name = "www"
	alternate_names = ["docs.example.com", "status.example.com"]
}
```

## Argument Reference

The following argument(s) are supported:

* `domain_id` - (Required) The domain to be issued the certificate for
* `name` - (Required) The certificate name; use `""` for the root domain. Wildcard names are supported.
* `alternate_names` - (Optional) The certificate alternate names
* `auto_renew` - (Required) True if the certificate should auto-renew
* `signature_algorithm` - (Optional) The signature algorithm to use for the certificate
* `timeouts` - (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))

## Attribute Reference

The following additional attributes are exported:

* `id` - The certificate ID
* `years` - The years the certificate will last
* `state` - The state of the certificate
* `authority_identifier` - The identifying certification authority (CA)
* `csr` - The certificate signing request
* `expires_at` - The datetime the certificate will expire
* `created_at` - The datetime the certificate was created
* `updated_at` - The datetime the certificate was last updated

<a id="nestedblock--timeouts"></a>

### Nested Schema for `timeouts`

Optional:

- `read` (String) - The timeout for the read operation
