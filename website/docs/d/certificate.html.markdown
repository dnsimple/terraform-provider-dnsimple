---
layout: "dnsimple"
page_title: "DNSimple: dnsimple_certificate"
sidebar_current: "docs-dnsimple-datasource-certificate"
description: |-
  Provides a DNSimple certificate data source.
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
