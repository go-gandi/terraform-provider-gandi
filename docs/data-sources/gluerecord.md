# Data: gluerecord

Use this data source to get the IPs listed within a Glue Record for a domain.

## Example Usage

```terraform
data "gandi_gluerecord" "my_domain" {
    zone = "my.domain"
    name = "ns1"
}
```

## Argument Reference

- `zone` - (Required) The domain name.
- `name` - (Required) The host name of the record.

## Attribute Reference

- `ips` - A list of the ip addresses provided in the glue record.

