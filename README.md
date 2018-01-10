# Terraform Gandi Provider

This provider allows managing DNS records on the Gandi LiveDNS service.

**Gandi is currently (as of Nov. 2017) migrating on a new platform, this provider is for the NEW platform.**

## Compiling

```
go get
go build -o terraform-provider-gandi
```

## Example

This example partly mimics the steps of [the official LiveDNS documentation example](http://doc.livedns.gandi.net/#quick-example), using the parts that have been implemented as Terraform resources.
Note: sharing_id is optional. It is used e.g. when the API key is registered to a user, where the domain you want to manage is not registered with that user (but the user does have rights on that zone/organization). 
```
provider "gandi" {
  key = "<the API key>"
  sharing_id = "<the sharing_id>"
}

resource "gandi_zone" "example_com" {
  name = "example.com Zone"
}

resource "gandi_zonerecord" "www" {
  zone = "${gandi_zone.example_com.id}"
  name = "www"
  type = "A"
  ttl = 3600
  values = [
    "192.168.0.1"
  ]
}

resource "gandi_domainattachment" "example_com" {
    domain = "example.com"
    zone = "${gandi_zone.example_com.id}"
}
```

This example sums up the available resources.
