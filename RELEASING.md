# Releasing

This document describes the steps to release a new version of the Terraform DNSimple Provider.

## Prerequisites

- You have commit access to the repository
- You have push access to the repository
- You have a GPG key configured for signing tags

## Release process

The following instructions use `$VERSION` as a placeholder, where `$VERSION` is a `MAJOR.MINOR.BUGFIX` release such as `1.2.0`.

1. **Run the test suite** and ensure all the tests pass

   ```shell
   make test
   ```

2. **Finalize the changelog** with the new version

   Edit `CHANGELOG.md` and finalize the `## main` section, assigning the version.

3. **Commit and push the changes**

   ```shell
   git commit -a -m "Release $VERSION"
   git push origin main
   ```

4. **Wait for CI to complete**

   Ensure the CI build passes on the main branch before proceeding.

5. **Create a signed tag**

   ```shell
   git tag -a v$VERSION -s -m "Release $VERSION"
   git push origin --tags
   ```

6. **CI and goreleaser will handle the rest**

   The CI workflow will automatically build and publish the release using goreleaser.

## Post-release

- Verify the new version appears on the [Terraform Registry](https://registry.terraform.io/providers/dnsimple/dnsimple)
