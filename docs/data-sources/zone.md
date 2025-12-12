---
page_title: "DNSimple: dnsimple_zone"
---

# dnsimple\_zone

Get information about a DNSimple zone.

It is generally preferable to use the `dnsimple_zone` resource, but you may wish to only retrieve and link the zone information when the resource exists in multiple Terraform projects.

## Example Usage

```hcl
data "dnsimple_zone" "example" {
  name = "example.com"
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Required) The name of the zone.

## Attributes Reference

The following attributes are exported:

- `id` - The zone ID.
- `account_id` - The account ID.
- `reverse` - Whether the zone is a reverse zone (`true`) or forward zone (`false`).
- `secondary` - Whether the zone is a secondary zone.
- `active` - Whether the zone is active.
- `last_transferred_at` - The last time the zone was transferred (only applicable for secondary zones).
