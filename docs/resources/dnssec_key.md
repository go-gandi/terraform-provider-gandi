# Resource: dnssec_key

The `dnssec_key` resource creates a DNSSEC key to a domain.

## Example Usage

```terraform
resource "gandi_dnssec_key" "my_key" {
    domain = "example.com"
    algorithm = 15
    type = "ksk"
    public_key = "Z6eCbfmpYPYmOJ0PYKq8fKzxcP3K/xEBlF5omvO+UwY="
}
```

## Argument Reference

- `domain` - (Required) The domain to add the key to.
- `name` - (Required) The algorithm used for the key.
- `type` - (Required) "ksk" or "zsk".
- `public_key` - (Required) The public key to use.

## Attribute Reference

- `digest` - Digest of the added key.
- `digest_type` - Type of digest.
- `keytag` - The keytag assigned by the server.

## Import

Existing keys can be imported by running:

`terraform import gandi_dnssec_key.my_key example.com/<id>`, where the key ID is a UUID that can be found through querying the [Gandi API](https://api.gandi.net/docs/domains/#get-v5-domain-domains-domain-dnskeys)
