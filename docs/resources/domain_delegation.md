---
page_title: "DNSimple: dnsimple_domain_delegation"
---

# dnsimple\_domain\_delegation

Provides a DNSimple domain delegation resource.

This resource works differently to other resources. It does not create a new resource, but instead updates the existing domain delegation. This is because the domain delegation is a property of the domain, not a separate resource. When this resource is destroyed, only the Terraform state is removed; the domain delegation is left intact and unmanaged by Terraform.

## Example Usage

```hcl
# Create a domain delegation
resource "dnsimple_domain_delegation" "foobar" {
  id = "${var.dnsimple.domain}"
  name_servers = ["example.org", "example.com"]
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