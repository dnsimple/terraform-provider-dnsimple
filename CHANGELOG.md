# Changelog

## 2.0.0 - 2025-12-15

BREAKING CHANGES:

- provider: Minimum Terraform version is now 1.12. Previous versions are no longer supported.

ENHANCEMENTS:

- provider: Add `debug_transport_file` configuration option to enable transport debug logging
- resource/`dnsimple_zone_record`: Validate format of record type (#317)
- deps: Bump dnsimple-go to v7


## 1.10.0 - 2025-07-08

ENHANCEMENTS:

- deps: Bump dnsimple-go to v5
- ci: Add terraform 1.12 to the test suite

NOTES:

- The `dnsimple_zone` data source is no longer deprecated and will not be removed in a future release.

## 1.9.1 - 2025-05-27

ENHANCEMENTS:

- deps: Bump golang.org/x/net from 0.37.0 to 0.38.0 (#269)
- make: Update .PHONY

NOTES:

- Delete CODEOWNERS
- Remove unnecessary gitignore entries

BUG FIXES:

- ci: Fixes CI / Acceptance Tests (#274)
- ci: Switch to Gofumpt (#270)
- ci: Format YAML
- ci: Align dependabot config

## 1.9.0 - 2025-04-09

BUG FIXES:

- deps: Upgrade `dnsimple-go` to v4.0.0 which ships a fix for (#428) (#264)
- resource/`dnsimple_domain_delegation`: Use plan data instead of state in domain delegation updates fixing (#256) (#266)
- resource/`dnsimple_zone_record`: Skip prefetch cache on resource import fixing (#238) (#267)

NOTES:

- We have updated the Go module to Go 1.24.

## 1.8.0 - 2024-10-17

NOTES:

- Updates the Terraform Plugin Framework to latest version (v1.10.0). In addition to other dependency updates.
- We have updated the Go module to Go 1.23.

## 1.7.0 - 2024-08-01

NOTES:

- Updates the Terraform Plugin Framework to latest version (v1.10.0). In addition to other dependency updates.

## 1.6.0 - 2024-05-28

ENHANCEMENTS:

- **Update Data Source:** `dnsimple_certificate` has been updated to have a stable ID. (#222)

## 1.5.0 - 2024-03-19

ENHANCEMENTS:

- **Update Resource:** `dnsimple_zone_record` has been updated to handle cases where the `name` attribute is normalized by the API, resulting in bad state as config differs from state.
- **Update Resource:** `dnsimple_domain_delegation` now has the trailing dot removed from the `name_servers` attribute entries. This is to align with the API and avoid perma diffs. (#203)

BUG FIXES:

- provider: Corrects the method by which the prefetch configuration flag is loaded from the environment. (#206)
- provider: Introduces concurrent read/write locking for the cache to prevent panics during simultaneous map writes. (#206)
- provider: Adjusts the logic for searching zone records in the cache, utilizing the normalized content value rather than the initially configured value. (#206)

NOTES:

- This release is no longer compatible with Terraform versions < 1.3. This is due to the new protocol changes in the underlying terraform framework. If you are using Terraform 1.3 or later, you should be unaffected by this change.
- We have updated the Go module to Go 1.21.

## 1.4.0 - 2024-01-17

FEATURES:

- **New Resource:** `dnsimple_zone` (#184)

NOTES:

- The `dnsimple_zone` data source is now deprecated and will be removed in a future release. Please migrate to the `dnsimple_zone` resource.
- This Go module has been updated to Go 1.20.

## 1.3.1 - 2023-11-02

BUG FIXES:

- resource/`dnsimple_zone_record`: As part of increasing our RFC compliance for TXT records we may normalize user-provided content and as a result unless the user provides a normalized version of the record value we will incorrectly write the normalized value to the state resulting in mismatch between the configuration and the written state. (#167)

## 1.3.0 - 2023-09-21

FEATURES:

- **New Data Source:** `dnsimple_registrant_change_check` (#155)
- **Update Resource:** `dnsimple_registered_domain` now supports the change of `contact_id` which results in domain contact change at the registry (#155)

BUG FIXES:

- resource/`dnsimple_registered_domain`: Boolean attributes could get toggled to `false` when the resource was updated. This applied for `auto_renew_enabled`, `whois_privacy_enabled`, `dnssec_enabled` and `transfer_lock_enabled`. This would happen if the attribute was not explicitly set in the configuration. Please check and update your configuration to explicitly specify the state of the attributes if you think you might be affected by this issue. (#155)

NOTES:

- The `contact_id` attribute previously only supported configuring at resource creation. Now the attribute can be changed and it will result in registrant change being initiated at the registrar. In addition the `extended_attributes` attribute can now also be updated after a domain has been registered and the values in the extended attributes will be passed along the registrant change as some TLDs require extended attributes to be passed along the registrant change.
- The registrant change can happen asynchronously for certain contact changes and TLDs. The resource supports this by attempting to sync the state of the registrant change in a similar way to how domain registration works with the `dnsimple_registered_domain` resource.

## 1.2.1 - 2023-09-19

ENHANCEMENTS:

- **Update Resource:** `dnsimple_zone_record` now supports the `regions` attribute which you can use to specify regional records (#156)

## 1.2.0 - 2023-09-08

FEATURES:

- **Update Resource:** `dnsimple_registered_domain` now supports `transfer_lock_enabled` argument which you can use to manage the domain transfer lock state of your registered domains (#143)

## 1.1.2 - 2023-05-17

BUG FIXES:

- resource/`dnsimple_lets_encrypt_certificate`: Fix error when no alternate names are provided for the Let's Encrypt certificate resource (#111)

## 1.1.1 - 2023-05-09

ENHANCEMENTS:

- resource/`dnsimple_registered_domain`: Support resource importing with domain name only (#107)

NOTES:

- Prior to this release the `dnsimple_registered_domain` resource could only be imported using the domain name and the domain registration ID. This release adds support for importing the resource using the domain name only. This has no effects on existing resources.

## 1.1.0 - 2023-04-21

FEATURES:

- **New Resource:** `dnsimple_contact` (#98)
- **New Resource:** `dnsimple_domain_delegation` (#99)
- **New Resource:** `dnsimple_ds_record` (#101)
- **New Resource:** `dnsimple_registered_domain` (#100)

ENHANCEMENTS:

- resource/`dnsimple_lets_encrypt_certificate`: Add `alternate_names` attribute (#102)

## 1.0.0 - 2023-04-13

NOTES:

- We've reached a stable 1.0.0 release! This is identical to version 0.17.0, but the API is now stable and we will follow semantic versioning from now on.
- If you are migrating from version 0.16.3 or earlier, refer to the changelog for 0.17.0 for the breaking changes.

## 0.17.0 - 2023-04-07

BREAKING CHANGES:

- provider: Drop support for Terraform 0.14 (#93)

- resource/`dnsimple_lets_encrypt_certificate`:
  - The deprecated `contact_id` field has been removed from the `dnsimple_lets_encrypt_certificate` resource. (#93)
  - The `id` field on the `dnsimple_lets_encrypt_certificate` is now of type `int64` instead of `string` to keep in line with the API. (#93)
  - The `domain_id` field on the `dnsimple_lets_encrypt_certificate` is now required. (#93)
  - The `expires_on` attribute on the `dnsimple_lets_encrypt_certificate` has been renamed to `expires_at` to keep in line with the API. (#93)

- resource/`dnsimple_zone_record`:
  - The `ttl` and `priority` fields on the `dnsimple_zone_record` are now of type `int64` instead of `string`. (#93)

- resource/`dnsimple_record`:
  - The resource has been removed from the provider as it was deprecated in v0.9.2. (#93)

- provider: The `PREFETCH` environment variable has been renamed to `DNSIMPLE_PREFETCH` to avoid conflicts with other services. (#93)

## 0.16.3 - 2023-03-23

BUG FIXES:

- resource/`dnsimple_zone_record`: Correctly error out and terminate import operations when invalid (#88)

## 0.16.2 - 2023-03-06

BUG FIXES:

- provider: Prefetch cache lookups are deterministic (#84)

## 0.16.1 - 2023-03-01

ENHANCEMENTS:

- docs: Improved the documentation for the `dnsimple_domain` and `dnsimple_lets_encrypt_certificate` resources. (#82)

## 0.16.0 - 2023-02-28

ENHANCEMENTS:

- resource/`dnsimple_lets_encrypt_certificate`: Added support for `signature_algorithm` attribute
- deps: Dependency updates

## 0.15.0 - 2022-11-22

NOTES:

- resource/`dnsimple_lets_encrypt_certificate`: Deprecate the use of `contact_id` in the `dnsimple_lets_encrypt_certificate` resource. The field is no longer in use and there is no replacement for it (#62)

ENHANCEMENTS:

- provider: Surface all API exceptions during a terraform run (#61)

BUG FIXES:

- resource/`dnsimple_zone_record`: Fixed error while importing record with underscore (#7)

## 0.14.1 - 2022-09-27

BUG FIXES:

- provider: Avoid panic when looking for a record and it does not exist on the prefetched list

## 0.14.0 - 2022-09-20

ENHANCEMENTS:

- provider: Pass parent context to DNSimple client calls to propagate errors and handling cancellation
- deps: Updated minimum go version to 1.18
- deps: Updated the `dnsimple-go` dependency to v1.0.0
- provider: Show validation errors when applying and point to the field which is failing

## 0.13.0 - 2022-06-17

ENHANCEMENTS:

- provider: Added ability to pass a custom user agent fragment (#56)

## 0.12.0 - 2022-06-15

ENHANCEMENTS:

- deps: Updated minimum go version to 1.17
- deps: Updated the terraform-plugin-sdk to v2.17.0
- provider: Set the token as sensitive so it is not logged

## 0.11.3 - 2022-04-13

ENHANCEMENTS:

- docs: Fixed documentation

## 0.11.2 - 2022-04-13

ENHANCEMENTS:

- docs: Added helpful links to the documentation

## 0.11.1 - 2022-02-15

ENHANCEMENTS:

- docs: Added the documentation for the `resource_dnsimple_lets_encrypt_certificate_resource`

## 0.11.0 - 2021-11-17

FEATURES:

- **New Data Source:** `dnsimple_certificate`
- **New Resource:** `dnsimple_lets_encrypt_certificate` to purchase and issue Let's Encrypt certificates

ENHANCEMENTS:

- resource/`dnsimple_domain`: Added the `dnsimple_domain` import to import domains
- deps: Updated the `dnsimple-go` dependency to v0.71.0

## 0.10.0 - 2021-10-22

ENHANCEMENTS:

- provider: Added the `prefetch` option to avoid running into API rate limitations when dealing with big configurations

## 0.9.2 - 2021-09-03

NOTES:

- resource/`dnsimple_record`: Added the deprecated `resource_dnsimple_record`

## 0.9.1 - 2021-09-03

NOTES:

- resource/`dnsimple_record`: Bring the `dnsimple_record` configuration back and adds a deprecation warning

## 0.9.0 - 2021-09-02

FEATURES:

- **New Data Source:** `dnsimple_zone`
- **New Resource:** `dnsimple_domain` (to create domains in DNSimple)

## 0.6.0 - 2021-07-22

NOTES:

- deps: Migrated SDK to v2 (version 2.7.0)
- deps: Updated dependencies

## 0.5.3 - 2021-05-14

NOTES:

- build: Include darwin_arm64 builds in release
- build: Removes /vendor directory

## 0.5.2 - 2021-05-14

NOTES:

- deps: Updated dependencies
- build: Include darwin_arm64 builds

## 0.5.1 - 2020-11-30

NOTES:

- ci: Move to GH Actions for publishing

## 0.5.0 - 2020-11-30

FEATURES:

- **New Resource:** `dnsimple_email_forward` ([#28](https://github.com/terraform-providers/terraform-provider-dnsimple/pull/28), [#30](https://github.com/terraform-providers/terraform-provider-dnsimple/pull/30))
- provider: Ability to switch to sandbox environment ([#12](https://github.com/terraform-providers/terraform-provider-dnsimple/pull/12))

## 0.4.0 - 2020-05-12

ENHANCEMENTS:

- deps: Upgraded to dnsimple-go v0.61.0

## 0.3.0 - 2020-02-11

NOTES:

- deps: Upgraded plugin to use the Terraform Plugin SDK v1.0.0 instead of Terraform Core ([#21](https://github.com/terraform-providers/terraform-provider-dnsimple/pulls/21))
- provider: Remove support for deprecated API v1 attributes ([#22](https://github.com/terraform-providers/terraform-provider-dnsimple/pulls/22))

ENHANCEMENTS:

- deps: Upgraded to dnsimple-go v0.31.0 ([#23](https://github.com/terraform-providers/terraform-provider-dnsimple/pulls/23))

## 0.2.0 - 2019-06-20

NOTES:

- This release includes a Terraform upgrade with compatibility for Terraform v0.12. The provider remains backwards compatible with Terraform v0.11 and there should not be any significant behavioural changes. ([#16](https://github.com/terraform-providers/terraform-provider-dnsimple/issues/16))

## 0.1.0 - 2017-06-20

NOTES:

- Same functionality as that of Terraform 0.9.8. Repacked as part of [Provider Splitout](https://www.hashicorp.com/blog/upcoming-provider-changes-in-terraform-0-10/)
