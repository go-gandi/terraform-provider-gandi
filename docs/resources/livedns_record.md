# Resource: livedns_record

The `livedns_record` resource creates a record for a domain in the LiveDNS management system.
!> You must enable a domain using the `livedns_domain` resource before adding records to it.

## Example Usage

```terraform
resource "gandi_livedns_record" "www_my_domain" {
    zone = "my.domain"
    name = "www"
    type = "A"
    values = ["127.0.0.2"]
}
```

## Argument Reference

* `zone` - (Required) The FQDN of the domain.
* `name` - (Required) The name of the record.
* `type` - (Required) The type of the record. Can be one of "A", "AAAA", "ALIAS", "CAA", "CDS", "CNAME", "DNAME", "DS", "KEY", "LOC", "MX", "NS", "OPENPGPKEY", "PTR", "SPF", "SRV", "SSHFP", "TLSA", "TXT", "WKS".
* `ttl` - (Required) The TTL of the record.
* `values` - (Required) A list of the values of the record.

## Attribute Reference

* `id` - The ID of the resource.
