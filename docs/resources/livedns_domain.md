# Resource: livedns_domain

The `livedns_domain` resource enables a Domain in the LiveDNS management system. You must enable a domain before adding records to it.

## Example Usage

```terraform
resource "gandi_livedns_domain" "my_domain" {
    name = "my.domain"
    ttl = "3600"
}
```

## Argument Reference

* `name` - (Required) The FQDN of the domain.
* `ttl` - (Required) The default TTL of the domain.
* `automatic_snapshots` - (Optional) Enable the automatic creation of snapshots when changes are made.

## Attribute Reference

* `id` - The ID of the resource.
