# Terraform Provider for DNSimple

[![Terraform Registry](https://img.shields.io/badge/terraform-registry-623CE4?logo=terraform)](https://registry.terraform.io/providers/dnsimple/dnsimple)
[![License](https://img.shields.io/badge/license-MPL--2.0-blue.svg)](LICENSE)

The Terraform DNSimple provider allows you to manage DNSimple resources using Terraform.

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.12
- [Go](https://golang.org/doc/install) >= 1.18 (to build the provider from source)

## Installation

The provider is available on the [Terraform Registry](https://registry.terraform.io/providers/dnsimple/dnsimple). Add the following to your Terraform configuration:

```hcl
terraform {
  required_providers {
    dnsimple = {
      source  = "dnsimple/dnsimple"
      version = "~> 1.0"
    }
  }
}

provider "dnsimple" {
  token   = var.dnsimple_token
  account = var.dnsimple_account
  sandbox = true  # Set to false for production
}
```

Then run:

```shell
terraform init
```

## Documentation

Full documentation is available on the [Terraform Registry](https://registry.terraform.io/providers/dnsimple/dnsimple/latest/docs).

## Quick Start

After installing the provider, configure it with your DNSimple credentials:

```hcl
provider "dnsimple" {
  token   = "your-api-token"
  account = "your-account-id"
  sandbox = true
}

resource "dnsimple_zone" "example" {
  name = "example.com"
}

resource "dnsimple_zone_record" "www" {
  zone_id = dnsimple_zone.example.id
  name    = "www"
  type    = "A"
  value   = "1.2.3.4"
}
```

## Development

### Getting Started

Clone the repository:

 ```shell
 git clone git@github.com:dnsimple/terraform-provider-dnsimple.git
 cd terraform-provider-dnsimple
 ```

Build the provider:

 ```shell
 make build
 ```

 This will build the provider and place the binary in `$GOPATH/bin`.

### Testing

Run the unit tests:

```shell
make test
```

Run the acceptance tests (requires DNSimple API credentials):

```shell
DNSIMPLE_ACCOUNT=12345 DNSIMPLE_TOKEN="your-token" DNSIMPLE_DOMAIN=example.com DNSIMPLE_SANDBOX=true make testacc
```

**Note:** Acceptance tests create real resources and may incur costs.

#### Testing Let's Encrypt Resources

The sandbox environment does not support certificate operations. To test `dnsimple_lets_encrypt_certificate` resources, run tests in production:

```shell
DNSIMPLE_SANDBOX=false DNSIMPLE_CERTIFICATE_NAME=www DNSIMPLE_CERTIFICATE_ID=123 make testacc
```

### Sideloading the Provider

To use a locally built version of the provider:

1. Install the provider:

   ```shell
   make install
   ```

2. Create a symlink to the Terraform plugins directory:

   ```shell
   # Replace darwin_arm64 with your architecture
   mkdir -p ~/.terraform.d/plugins/terraform.local/dnsimple/dnsimple/0.1.0/darwin_arm64
   ln -s "$GOBIN/terraform-provider-dnsimple" ~/.terraform.d/plugins/terraform.local/dnsimple/dnsimple/0.1.0/darwin_arm64/
   ```

3. Configure Terraform to use the local provider:

   ```hcl
   terraform {
    required_providers {
      dnsimple = {
        source  = "terraform.local/dnsimple/dnsimple"
        version = "0.1.0"
      }
    }
   }
   ```

4. Test with the example configuration:

   ```shell
   cd example
   terraform init && terraform apply
   ```

## Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines on how to contribute to this project.

## License

This project is licensed under the Mozilla Public License 2.0. See [LICENSE](LICENSE) for details.

## Resources

- [Terraform Registry](https://registry.terraform.io/providers/dnsimple/dnsimple)
- [DNSimple API Documentation](https://developer.dnsimple.com/)
- [DNSimple Support](https://support.dnsimple.com/)
