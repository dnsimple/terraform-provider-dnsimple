---
page_title: "DNSimple: dnsimple_zone"
---

# dnsimple\_zone

Provides a DNSimple zone resource.

~> **Note:** Currently the resource creation acts as an import, so the zone must already exist in DNSimple. The only attribute that will be modified during resource creation is the `active` state of the zone. This is because our API does not allow for the creation of zones. Creation of zones happens through the purchase or creation of domains. We expect this behavior to change in the future.

## Example Usage

```hcl
resource "dnsimple_zone" "example" {
  name = "example.com"
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Required) The zone name.

## Attributes Reference

- `id` - The ID of this resource.
- `account_id` - The account ID for the zone.
- `reverse` - Whether the zone is a reverse zone.
- `secondary` - Whether the zone is a secondary zone.
- `active` - Whether the zone is active.
- `last_transferred_at` - The last time the zone was transferred only applicable for **secondary** zones.

## Import

DNSimple zones can be imported using the zone name.

```bash
terraform import dnsimple_zone.example example.com
```

The zone name can be found within the [DNSimple Zones API](https://developer.dnsimple.com/v2/zones/#getZone). Check out [Authentication](https://developer.dnsimple.com/v2/#authentication) in API Overview for available options.
