# Terraform Gandi Provider

This provider supports [managing DNS zones](https://api.gandi.net/docs/domains/) and [managing the LiveDNS service](https://api.gandi.net/docs/livedns/) in Gandi.

This provider currently doesn't support the Email, Organization or Billing APIs. We welcome pull requests to implement more functionality!

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 0.12.x
- [Go](https://golang.org/doc/install) >= 1.12

## Installation

1. Clone the repository
1. Enter the repository directory
1. Build the provider:

```shell
make
make install
```

Once installed, run `terraform init` to enable the Gandi plugin in your terraform environment.

See the [Hashicorp Terraform documentation](https://www.terraform.io/docs/plugins/basics.html#installing-plugins) for further details.

## Using the provider

This example partly mimics the steps of [the official LiveDNS documentation example](http://doc.livedns.gandi.net/#quick-example), using the parts that have been implemented as Terraform resources.
Note: sharing_id is optional. It is used e.g. when the API key is registered to a user, where the domain you want to manage is not registered with that user (but the user does have rights on that zone/organization).

```terraform
terraform {
  required_providers {
    gandi = {
      versions = ["2.0.0-rc3"]
      source   = "github/go-gandi/gandi"
    }
  }
}

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

```terraform
terraform {
  required_providers {
    gandi = {
      versions = ["X.Y.Z"]
      source   = "github/go-gandi/gandi"
    }
  }
}

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

## Development

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `make`.

### Linting

We use [pre-commit](https://pre-commit.com/) to managing and maintaining hooks, you can follow the [official website instructions](https://pre-commit.com/#install) to install it.

**Install**

```bash
python3 -m pip install pre-commit
```

Then in the repo root dir

```bash
pre-commit install
```

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```shell
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.
