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
* `contact_id` - (Deprecated) The contact id for the certificate

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
