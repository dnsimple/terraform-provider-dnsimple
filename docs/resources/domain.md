---
page_title: "DNSimple: dnsimple_domain"
---

# dnsimple\_domain

Provides a DNSimple domain resource.

## Example Usage

```hcl
# Create a domain
resource "dnsimple_domain" "foobar" {
  name = "${var.dnsimple.domain}"
}
```

## Argument Reference

The following argument(s) are supported:

* `name` - (Required) The domain name to be created

## Import

DNSimple domains can be imported using their numeric record ID.

```
$ terraform import dnsimple_domain.resource_name 5678
```

The record ID can be found within [DNSimple Domains API](https://developer.dnsimple.com/v2/domains/#listDomains). Check out [Authentication](https://developer.dnsimple.com/v2/#authentication) in API Overview for available options.

```
$ curl -u 'EMAIL:PASSWORD' https://api.dnsimple.com/v2/1234/domains?name_like=example.com | jq
{
  "data": [
    {
      "id": 5678,
      "account_id": 1234,
      "registrant_id": null,
      "name": "example.com",
      "unicode_name": "example.com",
      "state": "hosted",
      "auto_renew": false,
      "private_whois": false,
      "expires_on": null,
      "expires_at": null,
      "created_at": "2021-10-01T00:00:00Z",
      "updated_at": "2021-10-01T00:00:00Z"
    }
  ],
  "pagination": {
    "current_page": 1,
    "per_page": 30,
    "total_entries": 1,
    "total_pages": 1
  }
}
```
