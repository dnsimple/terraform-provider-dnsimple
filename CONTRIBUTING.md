# Contributing to DNSimple/terraform-provider

## Getting started

#### 1. Clone the repository

Clone the repository and move into it:

```shell
git clone git@github.com:dnsimple/terraform-provider-dnsimple.git
cd terraform-provider-dnsimple
```

#### 2. Build and test

If you wish to work on the provider, you'll first need Go installed on your machine (version 1.18+ is required). You'll also need to correctly setup a GOPATH, as well as adding $GOPATH/bin to your $PATH.

To compile the provider, run make build. This will build the provider and put the provider binary in the $GOPATH/bin directory.

```shell
$ make build
...
$ $GOPATH/bin/terraform-provider-dnsimple
...
```


## Testing

```shell
make test
```

You can also run the integration tests like:

```shell
DNSIMPLE_ACCOUNT=12345 DNSIMPLE_TOKEN="adf23cf" DNSIMPLE_DOMAIN=example.com DNSIMPLE_SANDBOX=true make testacc
```

### Testing the let's encrypt resource and the certificate data-source

Our sandbox environment does not allow purchasing or issue certificates. For that reason, if you want to test the
`resource_dnsimple_lets_encrypt_certificate` you will have to run the tests in production
(setting `DNSIMPLE_SANDBOX=false` in the shell).

You will have to set the following env variables in your shell:
   - `DNSIMPLE_CERTIFICATE_NAME` the name for which to request the certificate i.e. **www**
   - `DNSIMPLE_CERTIFICATE_ID` the certificate ID used in the datasource test

## Sideload the plugin

Sideload the plugin

```shell
make install
# Replace `darwin_arm64` with your arch. GOBIN should be where the Go built binary is installed to.
ln -s "$GOBIN/terraform-provider-dnsimple" "$HOME/.terraform.d/plugins/terraform.local/dnsimple/dnsimple/0.1.0/darwin_arm64/."
```

Use this as the provider configuration:

```tf
dnsimple = {
  source  = "terraform.local/dnsimple/dnsimple"
  version = "0.1.0"
}
```

You can use the `./example/simple.tf` config to test the provider.

```shell
cd example
terraform init && terraform apply --auto-approve
```

## Releasing

The following instructions uses `$VERSION` as a placeholder, where `$VERSION` is a `MAJOR.MINOR.BUGFIX` release such as `1.2.0`.

1. Run the test suite and ensure all the tests pass.

1. Finalize the `## main` section in `CHANGELOG.md` assigning the version.

1. Commit and push the changes

    ```shell
    git commit -a -m "Release $VERSION"
    git push origin main
    ```

1. Wait for CI to complete.

1. Create a signed tag.

    ```shell
    git tag -a v$VERSION -s -m "Release $VERSION"
    git push origin --tags
    ```

1. CI and goreleaser will handle the rest
