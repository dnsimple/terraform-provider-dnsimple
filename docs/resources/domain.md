---
page_title: "DNSimple: dnsimple_domain"
---

# dnsimple\_domain

Provides a DNSimple domain resource.

!> **Warning** Destroying a `dnsimple_domain` resource deletes the domain from your DNSimple account, and **its DNS records are not recoverable**. A domain registered through DNSimple stays registered at the registry and remains yours until it expires, but recovering it means adding it back to your account as a zone and then [contacting DNSimple support](https://dnsimple.com/contact) to have the registration relinked. See [Recovering a deleted domain](https://support.dnsimple.com/articles/recovering-deleted-domain/).

Two things are worth knowing before you manage a domain with this resource:

- Set [`prevent_delete`](#prevent_delete) to `true` to make `terraform destroy` fail for the resource instead of deleting the domain.
- If you only want Terraform to stop tracking the domain, remove it from state rather than destroying it. See [Removing a domain from Terraform](#removing-a-domain-from-terraform).

## Example Usage

```hcl
resource "dnsimple_domain" "example" {
  name = "example.com"
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Required) The domain name to be created.
- `prevent_delete` - Optional flag to guard against [accidental registration deletion of the domain](https://support.dnsimple.com/articles/recovering-deleted-domain/) (default `false`). Set it to `true` and `terraform destroy` will fail for this resource. Because changing `name` requires replacement, and replacement destroys the existing domain, renaming also fails while it is enabled — set it back to `false` and apply first.

  ```hcl
  resource "dnsimple_domain" "example" {
    name           = "example.com"
    prevent_delete = true
  }
  ```

## Attributes Reference

- `id` - The ID of this resource.
- `account_id` - The account ID for the domain.
- `auto_renew` - Whether the domain is set to auto-renew.
- `private_whois` - Whether the domain has WhoIs privacy enabled.
- `trustee` - Whether the domain has a [trustee](https://support.dnsimple.com/articles/what-is-domain-trustee/) enabled.
- `registrant_id` - The ID of the registrant (contact) for the domain.
- `state` - The state of the domain.
- `unicode_name` - The domain name in Unicode format.

## Removing a domain from Terraform

To stop managing a domain with Terraform while leaving it in your DNSimple account, remove it from state instead of destroying it.

With a `removed` block (Terraform 1.7 and later), which lets the change go through a normal plan and apply:

```hcl
removed {
  from = dnsimple_domain.example

  lifecycle {
    destroy = false
  }
}
```

Or directly with the CLI:

```bash
terraform state rm dnsimple_domain.example
```

Both drop the resource from Terraform state and leave the domain untouched in your account. Remove the matching `resource` block from your configuration as part of the same change, otherwise the next plan proposes creating the domain again.

-> **Note** `prevent_delete` does not interfere with either approach. Neither one calls the DNSimple API, so a domain can be dropped from state while protection is still enabled.

## Import

DNSimple domains can be imported using the domain name.

```bash
terraform import dnsimple_domain.example example.com
```

The domain name can be found within the [DNSimple Domains API](https://developer.dnsimple.com/v2/domains/#listDomains). Check out [Authentication](https://developer.dnsimple.com/v2/#authentication) in API Overview for available options.
