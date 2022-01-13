# Gandi Provider

The Gandi provider enables the purchasing and management of the
following Gandi resources:

- [DNS zones](https://api.gandi.net/docs/domains/)
- [LiveDNS service](https://api.gandi.net/docs/livedns/)
- [Email](https://api.gandi.net/docs/email/)
- [SimpleHosting](https://api.gandi.net/docs/simplehosting/)

The provider needs to be configured with the proper credentials before it can be used.

Use the navigation to the left to read about the available resources.

## Example Usage

```terraform
terraform {
  required_providers {
    gandi = {
      version = "~> 2.0.0"
      source   = "go-gandi/gandi"
    }
  }
}

provider "gandi" {
  key = "MY_API_KEY"
}

resource "gandi_domain" "example_com" {
  name = "example.com"
}
```

## Authentication

The Gandi provider supports a couple of different methods for providing authentication credentials.

You can retrieve your API key by visiting the [Account Management](https://account.gandi.net/en/) screen, going to the `Security` tab and generating your `Production API Key`.

Optionally, you can provide a Sharing ID to specify an organization. If set, the Sharing ID indicates the organization that will pay for any ordered products, and will filter collections.

### Static Credentials

!> Hard-coding credentials into any Terraform configuration is not recommended, and risks leaking secrets should the configuration be committed to public version control.

Usage:

```terraform
provider "gandi" {
  key = "MY_API_KEY"
  sharing_id = "MY_SHARING_ID"
}
```

### Environment Variables

You can provide your credentials via the `GANDI_KEY` and `GANDI_SHARING_ID` environment variables, representing the API Key and the Sharing ID, respectively.

```terraform
provider "gandi" {}
```

Usage:

```terraform
$ export GANDI_KEY="MY_API_KEY"
$ export GANDI_SHARING_ID="MY_SHARING_ID"
$ terraform plan
```
