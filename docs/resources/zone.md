---
page_title: "DNSimple: dnsimple_zone"
---

# dnsimple\_zone

Provides a DNSimple zone resource.

-> Currently the resource creation acts as an import, so the zone must already exist in DNSimple. The only attribute that will be modified during resource creation is the `active` state of the zone. This is because our API does not allow for the creation of zones. Creation of zones happens through the purchase or creation of domains. We expect this behavior to change in the future.

## Example Usage

```hcl
# Create a zone
resource "dnsimple_zone" "foobar" {
  name = "${var.dnsimple.zone}"
}
```

## Argument Reference

The following argument(s) are supported:

- `name` - (Required) The zone name

## Attributes Reference

- `id` - The ID of this resource.
- `account_id` - The account ID for the zone.
- `reverse` - Whether the zone is a reverse zone.
- `secondary` - Whether the zone is a secondary zone.
- `active` - Whether the zone is active.
- `last_transferred_at` - The last time the zone was transferred only applicable for **secondary** zones.

## Import

DNSimple zones can be imported using their the zone name.

```bash
terraform import dnsimple_zone.resource_name example.com
```

The zone ID can be found within [DNSimple Zones API](https://developer.dnsimple.com/v2/zones/#getZone). Check out [Authentication](https://developer.dnsimple.com/v2/#authentication) in API Overview for available options.

```bash
curl -H 'Authorization: Bearer <ACCESS_TOKEN>' https://api.dnsimple.com/v2/1234/zones/example.com | jq
{
  "data": {
    "id": 1,
    "account_id": 1234,
    "name": "example.com",
    "reverse": false,
    "secondary": false,
    "last_transferred_at": null,
    "active": true,
    "created_at": "2023-04-18T04:58:01Z",
    "updated_at": "2024-01-16T15:53:18Z"
  }
}
```
