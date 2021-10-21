# Resource: domain

The Domain resource enables the creation and management of Domains.
!> Creating a new Domain will result in your account being charged.
~> It is not currently possible to delete Domains via the provider.

## Example Usage

```terraform
resource "gandi_domain" "my_domain" {
    name = "my.domain"
    autorenew = "true"
}
```

## Argument Reference

- `name` - (Required) The FQDN of the domain.
- `nameservers` - (Optional) A list of nameservers for the domain.
- `autorenew` - (Optional) Should the domain autorenew?
- `admin` - (Required) Nested block listing the admin details for the domain. See below for the structure of the contact detail blocks.
- `billing` - (Required) Nested block listing the billing details for the domain. See below for the structure of the contact detail blocks.
- `owner` - (Required) Nested block listing the owner details for the domain. See below for the structure of the contact detail blocks.
- `tech` - (Required) Nested block listing the admin details for the domain. See below for the structure of the contact detail blocks.

Nested contact detail blocks have the following structure:

- `country` - (Required) The two letter country code for the contact.
- `email` - (Required) The email address of the contact.
- `family_name` - (Required) The family name of the contact.
- `given_name` - (Required) The given name of the contact.
- `street_addr` - (Required) The street address of the contact.
- `type` - (Required) The type of contact. Can be one of `person`, `company`, `association` or `public`.
- `phone` - (Required) The phone number for the contact.
- `city` - (Required) The city the contact is based in.
- `state` - (Optional) The state code for the contact.
- `organisation` - (Required unless the `type` is `person`) The legal name of the organisation.
- `zip` - (Required) Postal Code/Zip of the contact.
- `data_obfuscated` - (Optional) Whether to obfuscate contact details in WHOIS. Defaults to `false`.
- `mail_obfuscated` - (Optional) Whether to obfuscate contact email in WHOIS. Defaults to `false`.
- `extra_parameters` - (Optional) Extra parameters, needed for some domain registries.

## Attribute Reference

- `id` - The ID of the resource.
