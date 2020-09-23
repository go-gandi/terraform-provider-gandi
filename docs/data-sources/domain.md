# Data Source: domain

Use this data source to get the ID of a domain for other resources.

## Example Usage

```terraform
data "gandi_domain" "my_domain" {
    name = "my.domain"
}
```

## Argument Reference

* `name` - (Required) The FQDN of the domain.

## Attribute Reference

* `id` - The ID of the domain.
* `nameservers` - A list of nameservers for the domain.
