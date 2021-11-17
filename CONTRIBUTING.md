# Contributing to DNSimple/terraform-provider

## Getting started

#### 1. Clone the repository

Clone the repository and move into it:

```shell
git clone git@github.com:dnsimple/terraform-provider-dnsimple.git
cd terraform-provider-dnsimple
```

#### 2. Build and test

If you wish to work on the provider, you'll first need Go installed on your machine (version 1.12+ is required). You'll also need to correctly setup a GOPATH, as well as adding $GOPATH/bin to your $PATH.

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

### Testing the let's encrypt resource

Our sandbox environment does not allow purchasing or issue certificates. For that reason, if you want to test the 
`resource_dnsimple_lets_encrypt_certificate` you will have to run the tests in production 
(setting `DNSIMPLE_SANDBOX=false` in the shell).

First you will have to go to the `resource_dnsimple_lets_encrypt_certificate_test` and change the `domain` (line 21) 
to a real domain ID you want test against.

After that you will have to change the `testAccLetsEncrypConfig` (in that same file) changing the arguments marked:
   - contact_id (required)
   - and name (optional, but you might have to change it if you run the tests for a second time)


## Sideload the plugin

Sideload the plugin

```shell
make install
```

You can use the `./example/simple.tf` config to test the provider.

```shell
cd example
terraform init && terraform apply --auto-approve
```

## Releasing

The following instructions uses `$VERSION` as a placeholder, where `$VERSION` is a `MAJOR.MINOR.BUGFIX` release such as `1.2.0`.

1. Run the test suite and ensure all the tests pass.

1. Finalize the `## master` section in `CHANGELOG.md` assigning the version.

1. Commit and push the changes

    ```shell
    git commit -a -m "Release $VERSION"
    git push origin master
    ```

1. Wait for CI to complete.

1. Create a signed tag.

    ```shell
    git tag -a v$VERSION -s -m "Release $VERSION"
    git push origin --tags
    ```

1. CI and goreleaser will handle the rest
