---
page_title: "DNSimple: dnsimple_zone"
---

# dnsimple\_zone

Get information about a DNSimple zone. It is generally preferrable to use the `dnsimple_zone` resource but you may wish to only retrieve and link the zone information when the resource exists in multiple Terraform projects.

# Example Usage

Get zone:

```hcl
data "dnsimple_zone" "foobar" {
    name = "dnsimple.com"
}
```

# Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the zone

# Attributes Reference

The following additional attributes are exported:

* `id` - The zone ID
* `account_id` - The account ID
* `reverse` - True for a reverse zone, false for a forward zone.
