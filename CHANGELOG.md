# Terraform Gandi provider changelog

## v2.1.0

### Added

- Added the `mutable` attribute on the `livedns_records`
  ressource. When this attribute is set to `true`, some elements in
  the `TXT` value list can be managed outside from Terraform. This
  allows bots to add entries to a `TXT` record managed by the Gandi
  Terraform provider: they won't be removed on the next `terraform apply`!

### Fixed

- When a record has been manually removed (from the web UI for
  instance), the provider now recreates it.
- When importing a domain resource, the provider no longer imports
  nameservers when if LiveDNS is enabled on the domain: LiveDNS
  manages the nameservers and they can be updated by Gandi without
  any user intervention, leading to Terraform state incoherences.

## v1.0.0

After months of being used with success, here is the first release!
