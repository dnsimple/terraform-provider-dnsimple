# Contributing to DNSimple/terraform-provider

## Getting started

#### 1. Clone the repository

Clone the repository and move into it:

```shell
git clone git@github.com:dnsimple/terraform-provider-dnsimple.git
cd terraform-provider-dnsimple
```

#### 2. Install Go & Terraform

#### 3. Install the dependencies

```shell
go get ./...
```

#### 4. Build and test

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

## Sideload the plugin

Sideload the plugin

```shell
make install
```

You can use the `./example/simple.tf` config to test the provider.

```
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
