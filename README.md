# <a href="https://github.com/go-gandi/"><img align="right" src="./go-gandi.svg" alt="go-gandi community logo" title="go-gandi community" style="width: 120px; height: 84px;" width="120" height="84"/></a> Terraform Gandi Provider

This provider supports managing the following [Gandi
resources](https://api.gandi.net/docs/):

- [DNS zones](https://api.gandi.net/docs/domains/)
- [LiveDNS service](https://api.gandi.net/docs/livedns/)
- [Email](https://api.gandi.net/docs/email/)
- [SimpleHosting](https://api.gandi.net/docs/simplehosting/)

This provider currently doesn't support Organization or Billing APIs. We welcome pull requests to implement more functionality!

This is not an officially supported Gandi project.

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

```terraform
terraform {
  required_providers {
    gandi = {
      source = "go-gandi/gandi"
      version = ">= 2.1.0"
    }
  }
}

provider "gandi" {
  personal_access_token = "<the Personal Access Token>"
}

resource "gandi_domain" "example_com" {
  name = "example.com"
  owner {
    email = "gandi@example.com"
    type = "person"
    street_addr = "Example"
    zip = "75000"
    phone = "+33666666666"
    given_name = "Gandi"
    family_name = "Net"
    country = "FR"
    city = "Paris"
  }
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

resource "gandi_gluerecord" "example_com_gluerecord" {
    zone = "${gandi_livedns_domain.example_com.id}"
    name = "ns1"
    ips = [
        "1.1.1.1"
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
      version = "X.Y.Z"
      source  = "github/go-gandi/gandi"
    }
  }
}

provider "gandi" {
  personal_access_token = "<the Personal Access Token>"
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

This provider is distributed under the terms of the Mozilla Public
License version 2.0 (see the [`LICENSE` file](./LICENSE)).

Its initial author (@tiramiseb) is not affiliated in any way with
Gandi - apart from being a happy customer of their services.

## Development

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `make`.

To use the provider, run `make install` to populate your local
Terraform provider cache
(`$(HOME)/.terraform.d/plugins/registry.terraform.io`). On a next
`terraform init`, Terraform then picks the provider from this directory.

### Linting

We use [pre-commit](https://pre-commit.com/) to manage and maintain the pre-commit hooks, you can follow the [official instructions](https://pre-commit.com/#install) to install it.

**Install**

```bash
python3 -m pip install pre-commit
```

Then in the repo root dir

```bash
pre-commit install
```

### Generating documentation

As suggested by the [Terraform
documentation](https://www.terraform.io/docs/registry/providers/docs.html#generating-documentation),
we use the
[tfplugindocs](https://github.com/hashicorp/terraform-plugin-docs)
project to generate the provider documentation.

To generate the documentation, just run:

     $ tfplugindocs

### Releasing

There are GitHub actions configured to generate a release from a tag
such as described in the [Terraform provider
documentation](https://learn.hashicorp.com/tutorials/terraform/provider-release-publish?in=terraform/providers)

To publish a new release:

    $ git tag vX.X
    $ git push origin vX.X

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```shell
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.
