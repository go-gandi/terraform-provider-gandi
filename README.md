# Terraform Gandi Provider

This provider supports [managing DNS zones](https://api.gandi.net/docs/domains/) and [managing the LiveDNS service](https://api.gandi.net/docs/livedns/) in Gandi.

This provider currently doesn't support the Email, Organization or Billing APIs. We welcome pull requests to implement more functionality!

## Installation

```
make
mkdir -p ~/.terraform.d/plugins/
install -m 644 terraform-provider-gandi ~/.terraform.d/plugins/
```

Once installed, run `terraform init` to enable the Gandi plugin in your terraform environment.

See the [Hashicorp Terraform documentation](https://www.terraform.io/docs/plugins/basics.html#installing-plugins) for further details.

## Example

This example partly mimics the steps of [the official LiveDNS documentation example](http://doc.livedns.gandi.net/#quick-example), using the parts that have been implemented as Terraform resources.
Note: sharing_id is optional. It is used e.g. when the API key is registered to a user, where the domain you want to manage is not registered with that user (but the user does have rights on that zone/organization). 

```
provider "gandi" {
  key = "<the API key>"
  sharing_id = "<the sharing_id>"
}

resource "gandi_domain" "example_com" {
  name = "example.com"
  nameservers = gandi_livedns_domain.example_com.nameservers
}

resource "gandi_livedns_domain" "example_com" {
  name = "example.com"
}

resource "gandi_livedns_record" "www_example_com" {
  zone = "${gandi_livedns_domain.example_com.id}"
  name = "www"
  type = "A"
  ttl = 3600
  values = [
    "192.168.0.1"
  ]
}
```

This example sums up the available resources.

### Zone data source

If your zone already exists (which is very likely), you may use it as a data source:

```
provider "gandi" {
  key = "<the API key>"
  sharing_id = "<the sharing_id>"
}

data "gandi_domain" "example_com" {
  name = "example.com"
}

resource "gandi_livedns_record" "www" {
  zone = "${data.gandi_domain.example_com.id}"
  name = "www"
  type = "A"
  ttl = 3600
  values = [
    "192.168.0.1"
  ]
}
```

## Licensing

This provider is distributed under the terms of the Mozilla Public License version 2.0. See the `LICENSE` file.

Its main author is not affiliated in any way with Gandi - apart from being a happy customer of their services.
