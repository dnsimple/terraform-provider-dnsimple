---
page_title: "DNSimple: dnsimple_certificate"
---

# dnsimple\_certificate

Get information about a DNSimple SSL certificate.

## Example Usage

```hcl
data "dnsimple_certificate" "example" {
  domain         = "example.com"
  certificate_id = 1234
}
```

## Argument Reference

The following arguments are supported:

- `domain` - (Required) The domain name of the SSL certificate.
- `certificate_id` - (Required) The ID of the SSL certificate.
- `timeouts` - (Block, Optional) (see [below for nested schema](#nested-schema-for-timeouts))

## Attributes Reference

The following attributes are exported:

- `id` - The certificate ID.
- `server_certificate` - The SSL certificate.
- `root_certificate` - The root certificate of the issuing CA.
- `certificate_chain` - A list of certificates that make up the certificate chain.
- `private_key` - The corresponding private key for the SSL certificate.

### Nested Schema for `timeouts`

Optional:

- `read` (String) - The timeout for the read operation, e.g., `5m`.
