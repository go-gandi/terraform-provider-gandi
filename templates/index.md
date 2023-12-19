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
  personal_access_token = "MY_PERSONAL_ACCESS_TOKEN"
}

resource "gandi_domain" "example_com" {
  name = "example.com"
}
```

## Authentication

The Gandi provider supports a couple of different methods for providing authentication credentials.

The recommended way is to create a Personal Access Token. Read more about these tokens in the [Gandi public API documentation](https://api.gandi.net/docs/authentication/).

The previous method of using an API key is now deprecated and should not be used anymore, though it is still supported by this provider for now. When using an API Key, you could also provide a Sharing ID to specify an organization. If set, the Sharing ID indicates the organization that will pay for any ordered products, and will filter collections.

### Static Credentials

!> Hard-coding credentials into any Terraform configuration is not recommended, and risks leaking secrets should the configuration be committed to public version control.

Usage:

```terraform
provider "gandi" {
  personal_access_token = "MY_PERSONAL_ACCESS_TOKEN"
}
```

### Environment Variables

You can provide your credentials via the `GANDI_PERSONAL_ACCESS_TOKEN` environment variable, representing the Personal Access Token.

```terraform
provider "gandi" {}
```

Usage:

```terraform
$ export GANDI_PERSONAL_ACCESS_TOKEN="MY_PERSONAL_ACCESS_TOKEN"
$ terraform plan
```
