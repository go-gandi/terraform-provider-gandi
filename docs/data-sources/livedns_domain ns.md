# Data Source: livedns_domain_ns

Use this data source to get the nameservers of a domain for other resources.

## Example Usage

```terraform
data "gandi_livedns_domain_ns" "my_domain" {
    name = "my.domain"
}
```

## Argument Reference

* `name` - (Required) The FQDN of the domain.

## Attribute Reference

* `id` - The ID of the domain.
* `nameservers` - A list of nameservers for the domain.
