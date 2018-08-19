# Terraform Gandi Provider

This provider allows managing DNS records on the Gandi LiveDNS service.

https://doc.livedns.gandi.net/

**This provider is only for the LiveDNS service: it does not deal with Gandi's other API, which allows managing other services**

## Is this provider alive?

I know we are generally worried when using libs with few commits and not much activity. Please note I am still using this provider (dogfooding, yay!), it is working and, I *will* fix problems if I encounter some and I *will* answer issues whenever needed (even if - shame on me - I sometimes take long time fixing what I consider minor issues).

(written on 2018-08-19)

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
