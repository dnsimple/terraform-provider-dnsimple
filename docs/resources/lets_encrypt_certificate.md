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
}
```

## Argument Reference

The following argument(s) are supported:

* `domain_id` - (Required) The domain to be issued the certificate for
* `name` - (Required) The certificate name
* `auto_renew` - (Required) True if the certificate should auto-renew
* `contact_id` - (Deprecated) The contact id for the certificate
* `signature_algorithm` - (Optional) The signature algorithm to use for the certificate
* `timeouts` - (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))

## Attribute Reference

The following attributes are exported:

* `id` - The certificate ID
* `domain_id` - The domain ID
* `contact_id` - The contact ID
* `name` - The certificate name
* `years` - The years the certificate will last
* `state` - The state of the certificate
* `authority_identifier` - The identifying certification authority (CA)
* `auto_renew` - Set to true if the certificate will auto-renew
* `csr` - The certificate signing request
* `expires_on` - The date the certificate will expire
* `created_at` - The datetime the certificate was created
* `updated_at` - The datetime the certificate was last updated

<a id="nestedblock--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `read` (String) - The timeout for the read operation
