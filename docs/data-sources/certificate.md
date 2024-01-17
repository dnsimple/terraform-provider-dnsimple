---
page_title: "DNSimple: dnsimple_certificate"
---

# dnsimple\_certificate

Provides a DNSimple certificate data source.

## Example Usage

```hcl
data "dnsimple_certificate" "foobar" {
  domain           = "${var.dnsimple_domain}"
  certificate_id   = "${var.dnsimple_certificate_id}"
}
```

## Argument Reference

The following arguments are supported:

* `domain` - (Required) The domain of the SSL Certificate
* `certificate_id` - (Required) The ID of the SSL Certificate

## Attributes Reference

The following attributes are exported:

* `server_certificate` - The SSL Certificate
* `root_certificate` - The Root Certificate of the issuing CA
* `certificate_chain` - A list of certificates that make up the chain
* `private_key` - The corresponding Private Key for the SSL Certificate

<a id="nestedblock--timeouts"></a>

### Nested Schema for `timeouts`

Optional:

- `read` (String) - The timeout for the read operation e.g. `5m`
