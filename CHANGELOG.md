# Changelog

## main

## 1.2.0 (Unreleased)

FEATURES:

- **New Resource:** `dnsimple_registered_domain_contact` (dnsimple/terraform-provider-dnsimple#142)
- **New Data Source:** `dnsimple_registrant_change_check` (dnsimple/terraform-provider-dnsimple#142)

## 1.1.2

BUG FIXES:

- Fix error when no alternate names are provided for the Let's Encrypt certificate resource (dnsimple/terraform-provider-dnsimple#111)

## 1.1.1

IMPROVEMENTS:

- `dnsimple_registered_domain`: Support resource importing with domain name only (dnsimple/terraform-provider-dnsimple#107)

NOTES:

Prior to this release the `dnsimple_registered_domain` resource could only be imported using the domain name and the domain registration ID. This release adds support for importing the resource using the domain name only. This has no effects on existing resources.

## 1.1.0

FEATURES:

- **New Resource:** `dnsimple_contact` (dnsimple/terraform-provider-dnsimple#98)
- **New Resource:** `dnsimple_domain_delegation` (dnsimple/terraform-provider-dnsimple#99)
- **New Resource:** `dnsimple_ds_record` (dnsimple/terraform-provider-dnsimple#101)
- **New Resource:** `dnsimple_registered_domain` (dnsimple/terraform-provider-dnsimple#100)

IMPROVEMENTS:

- `dnsimple_lets_encrypt_certificate`: Add `alternate_names` attribute (dnsimple/terraform-provider-dnsimple#102)

## 1.0.0

We've reached a stable 1.0.0 release! This is identical to version 0.17.0, but the API is now stable and we will follow semantic versioning from now on.

If you are migrating from version 0.16.3 or earlier, refer to the changelog for [0.17.0](#0170) for the breaking changes.

## 0.17.0

BREAKING CHANGES:

- Drop support for Terraform 0.14 (dnsimple/terraform-provider-dnsimple#93)

- Resource `dnsimple_lets_encrypt_certificate`:
  - The deprecated `contact_id` field has been removed from the `dnsimple_lets_encrypt_certificate` resource. (dnsimple/terraform-provider-dnsimple#93)
  - The `id` field on the `dnsimple_lets_encrypt_certificate` is now of type `int64` instead of `string` to keep in line with the API. (dnsimple/terraform-provider-dnsimple#93)
  - The `domain_id` field on the `dnsimple_lets_encrypt_certificate` is now required. (dnsimple/terraform-provider-dnsimple#93)
  - The `expires_on` attribute on the `dnsimple_lets_encrypt_certificate` has been renamed to `expires_at` to keep in line with the API. (dnsimple/terraform-provider-dnsimple#93)

- Resource `dnsimple_zone_record`:
  - The `ttl` and `priority` fields on the `dnsimple_zone_record` are now of type `int64` instead of `string`. (dnsimple/terraform-provider-dnsimple#93)

- Resource `dnsimple_record`:
  - The resource has been removed from the provider as it was deprecated in v0.9.2. (dnsimple/terraform-provider-dnsimple#93)

- The `PREFETCH` environment variable has been renamed to `DNSIMPLE_PREFETCH` to avoid conflicts with other services. (dnsimple/terraform-provider-dnsimple#93)

## 0.16.3

BUG FIXES:

* Correctly error out and terminate `dnsimple_zone_record` import operations when invalid (dnsimple/terraform-provider-dnsimple#88)

## 0.16.2

BUG FIXES:

* Prefetch cache lookups are deterministic (dnsimple/terraform-provider-dnsimple#84)

## 0.16.1

IMPROVEMENTS:

* Improved the documentation for the `dnsimple_domain` and `dnsimple_lets_encrypt_certificate` resources. (dnsimple/terraform-provider-dnsimple#82)

## 0.16.0

* Added support for `signature_algorithm` attribute in the `dnsimple_lets_encrypt_certificate` resource.
* Dependency updates.

## 0.15.0

* Deprecate the use of `contact_id` in the `dnsimple_lets_encrypt_certificate` resource. The field is no longer in use and there is no replacement for it (dnsimple/terraform-provider-dnsimple#62)
* Surface all API exceptions during a terraform run (dnsimple/terraform-provider-dnsimple#61)
* Fixed error while importing record with underscore (dnsimple/terraform-provider-dnsimple#7)

## 0.14.1

* Avoid panic when looking for a record and it does not exist on the prefetched list

## 0.14.0

* Pass parent context to DNSimple client calls to propagate errors and handling cancellation
* Updated minimum go version to 1.18
* Updated the `dnsimple-go` dependency to v1.0.0
* Show validation errors when applying and point to the field which is failing

## 0.13.0

* Added ability to pass a custom user agent fragment (dnsimple/terraform-provider-dnsimple#56)

## 0.12.0

* Updated minimum go version to 1.17
* Updated the terraform-plugin-sdk to v2.17.0
* Set the token as sensitive so it is not logged

## 0.11.3

* Fixed documentation

## 0.11.2

* Added helpful links to the documentation

## 0.11.1

* Added the documentation for the `resource_dnsimple_lets_encrypt_certificate_resource`.

## 0.11.0

* Added the `dnsimple_certificate` data-source.
* Added the `dnsimple_domain` import to import domains.
* Added the `resource_dnsimple_lets_encrypt_certificate_resource` to purchase and issue Let's Encrypt certificates.
* Updated the `dnsimple-go` dependency to v0.71.0

## 0.10.0

* Added the `prefetch` option to avoid running into API rate limitations when dealing with big configurations.

## 0.9.2

* Added the deprecated `resource_dnsimple_record`

## 0.9.1

* Bring the `dnsimple_record` configuration back and adds a deprecation warning

## 0.9.0

* Added the zone data-source
* Added the domain resource (to create domains in DNSimple)

## 0.6.0

NOTES

* Migrated SDK to v2 (version 2.7.0)
* Updated dependencies

## 0.5.3

NOTES

* Include darwin_arm64 builds in release

## 0.5.3

NOTES

* Removes /vendor directory

## 0.5.2

NOTES

* Updated dependencies
* Include darwin_arm64 builds

## 0.5.1

NOTES

* Move to GH Actions for publishing

## 0.5.0

FEATURES:

* **New Resource:** `dnsimple_email_forward` ([#28](https://github.com/terraform-providers/terraform-provider-dnsimple/pull/28), [#30](https://github.com/terraform-providers/terraform-provider-dnsimple/pull/30))
* Abilty to switch to sandbox environment ([#12](https://github.com/terraform-providers/terraform-provider-dnsimple/pull/12))

## 0.4.0 (May 12, 2020)

ENHANCEMENTS

* Upgraded to dnsimple-go v0.61.0

## 0.3.0 (February 11, 2020)

NOTES

* Upgraded plugin to use the Terraform Plugin SDK v1.0.0 instead of Terraform Core ([#21](https://github.com/terraform-providers/terraform-provider-dnsimple/pulls/21))
* Remove support for deprecated API v1 attributes ([#22](https://github.com/terraform-providers/terraform-provider-dnsimple/pulls/22))

ENHANCEMENTS

* Upgraded to dnsimple-go v0.31.0 ([#23](https://github.com/terraform-providers/terraform-provider-dnsimple/pulls/23))

## 0.2.0 (June 20, 2019)

NOTES

* This release includes a Terraform upgrade with compatibility for Terraform v0.12. The provider remains backwards compatible with Terraform v0.11 and there should not be any significant behavioural changes. ([#16](https://github.com/terraform-providers/terraform-provider-dnsimple/issues/16))

## 0.1.0 (June 20, 2017)

NOTES

* Same functionality as that of Terraform 0.9.8. Repacked as part of [Provider Splitout](https://www.hashicorp.com/blog/upcoming-provider-changes-in-terraform-0-10/)
