---
page_title: "DNSimple: dnsimple_zone"
---

# dnsimple\_zone

Get information about a DNSimple zone.

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

The following attributes are exported:

* `id` - The zone ID
* `account_id` - The account ID
* `name` - The name of the zone
* `reverse` - True for a reverse zone, false for a forward zone.
