---
page_title: "DNSimple: dnsimple_domain_delegation"
---

# dnsimple\_domain\_delegation

Provides a DNSimple domain delegation resource.

This resource allows you to control the delegation records (name servers) for a domain.

~> **Note:** This resource currently only supports the management of domains that are registered with DNSimple.
-> **Note:** When this resource is destroyed, only the Terraform state is removed; the domain delegation is left intact and unmanaged by Terraform.

## Example Usage

```hcl
# Create a domain delegation
resource "dnsimple_domain_delegation" "foobar" {
  id = "${var.dnsimple.domain}"
  name_servers = ["ns1.example.org", "ns2.example.com"]
}
```

## Argument Reference

The following argument(s) are supported:

* `id` - (Required) The domain ID.
* `name_servers` - (Required) The list of name servers to delegate to.

# Attributes Reference

There are no additional attributes.

## Import

DNSimple domain delegations can be imported.

```bash
terraform import dnsimple_domain.resource_name domain_id
```
