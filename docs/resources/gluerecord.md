# Resource: livedns_record

The `gluerecord` resource creates a Glue Record for a domain in the LiveDNS management system.

## Example Usage

```terraform
resource "gandi_gluerecord" "my_domain_gluerecord" {
    zone = "my.domain"
    name = "ns1"
    values = ["127.0.0.2"]
}
```

## Argument Reference

- `zone` - (Required) The FQDN of the domain.
- `name` - (Required) The name of the record.
- `values` - (Required) A list of the ip addresses of the record.

## Attribute Reference

- `id` - The ID of the resource.
