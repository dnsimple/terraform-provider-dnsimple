# Changelog

## main

## 1.10.0

ENHANCEMENTS:

- deps: Bump dnsimple-go to v5
- deps: Bump github.com/hashicorp/terraform-plugin-docs to 0.22.0
- deps: Bump github.com/hashicorp/terraform-plugin-testing to 1.13.2
- deps: Bump github.com/hashicorp/terraform-plugin-go to 0.28.0
- deps: Bump github.com/hashicorp/terraform-plugin-framework to 1.15.0
- deps: Bump github.com/cloudflare/circl to 1.6.1
- ci: Add terraform 1.12 to the test suite

NOTES:

- The `dnsimple_zone` data source is no longer deprecated and will not be removed in a future release.

## 1.9.1

ENHANCEMENTS:

- Bump golang.org/x/net from 0.37.0 to 0.38.0 (#269)
- make: Update .PHONY

NOTES:

- Delete CODEOWNERS
- Remove unnecessary gitignore entries

BUG FIXES:

- ci:
  - Fixes CI / Acceptance Tests (#274)
  - ci: Switch to Gofumpt (#270)
  - ci: Format YAML
  - ci: Align dependabot config

## 1.9.0

BUG FIXES:

- Upgrade `dnsimple-go` to v4.0.0 which ships a fix for (#428) (dnsimple/terraform-provider-dnsimple#264)
- Use plan data instead of state in domain delegation updates fixing (#256) (dnsimple/terraform-provider-dnsimple#266)
- Skip prefetch cache on `dnsimple_zone_record` resource import fixing (#238) (dnsimple/terraform-provider-dnsimple#267)

NOTES:

- We have updated the Go module to Go 1.24.

## 1.8.0

NOTES:

- Updates the Terraform Plugin Framework to latest version (v1.10.0). In addition to other dependency updates.
- We have updated the Go module to Go 1.23.

## 1.7.0

NOTES:

- Updates the Terraform Plugin Framework to latest version (v1.10.0). In addition to other dependency updates.

## 1.6.0

ENHANCEMENTS:

- **Update Data Source:** `dnsimple_certificate` has been updated to have a stable ID. (dnsimple/terraform-provider-dnsimple#222)

## 1.5.0

ENHANCEMENTS:

- **Update Resource:** `dnsimple_zone_record` has been updated to handle cases where the `name` attribute is normalized by the API, resulting in bad state as config differs from state.
- **Update Resource:** `dnsimple_domain_delegation` now has the trailing dot removed from the `name_servers` attribute entries. This is to align with the API and avoid perma diffs. (dnsimple/terraform-provider-dnsimple#203)

BUG FIXES:

- Corrects the method by which the prefetch configuration flag is loaded from the environment. (#206)
- Introduces concurrent read/write locking for the cache to prevent panics during simultaneous map writes. (#206)
- Adjusts the logic for searching zone records in the cache, utilizing the normalized content value rather than the initially configured value. (#206)

NOTES:

- This release is no longer compatible with Terraform versions < 1.3. This is due to the new protocol changes in the underlying terraform framework. If you are using Terraform 1.3 or later, you should be unaffected by this change.
- We have updated the Go module to Go 1.21.

## 1.4.0

FEATURES:

- **New Resource:** `dnsimple_zone` (dnsimple/terraform-provider-dnsimple#184)

NOTES:

- The `dnsimple_zone` data source is now deprecated and will be removed in a future release. Please migrate to the `dnsimple_zone` resource.
- This Go module has been updated to Go 1.20.

## 1.3.1

BUG FIXES:

- As part of increasing our RFC compliance for TXT records we may normalize user-provided content and as a result unless the user provides a normalized version of the record value we will incorrectly write the normalized value to the state resulting in mismatch between the configuration and the written state. (dnsimple/terraform-provider-dnsimple#167)

## 1.3.0

FEATURES:

- **New Data Source:** `dnsimple_registrant_change_check` (dnsimple/terraform-provider-dnsimple#155)
- **Updated Resource:** `dnsimple_registered_domain` now supports the change of `contact_id` which results in domain contact change at the registry (dnsimple/terraform-provider-dnsimple#155)

BUG FIXES:

- Boolean attributes in the `dnsimple_registered_domain` resource could get toggled to `false` when the resource was updated. This applied for `auto_renew_enabled`, `whois_privacy_enabled`, `dnssec_enabled` and `transfer_lock_enabled`. This would happen if the attribute was not explicitly set in the configuration. Please check and update your configuration to explicitly specify the state of the attributes if you think you might be affected by this issue. (dnsimple/terraform-provider-dnsimple#155)

NOTES:

The `contact_id` attribute previously only supported configuring at resource creation. Now the attribute can be changed and it will result in registrant change being initiated at the registrar. In addition the `extended_attributes` attribute can now also be updated after a domain has been registered and the values in the extended attributes will be passed along the registrant change as some TLDs require extended attributes to be passed along the registrant change.
The registrant change can happen asynchronously for certain contact changes and TLDs. The resource supports this by attempting to sync the state of the registrant change in a similar way to how domain registration works with the `dnsimple_registered_domain` resource.

## 1.2.1

ENHANCEMENTS:

- **Updated Resource:** `dnsimple_zone_record` now supports the `regions` attribute which you can use to specify regional records (dnsimple/terraform-provider-dnsimple#156)

## 1.2.0

FEATURES:

- **Updated Resource:** `dnsimple_registered_domain` now supports `transfer_lock_enabled` argument which you can use to manage the domain transfer lock state of your registered domains (dnsimple/terraform-provider-dnsimple#143)

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
