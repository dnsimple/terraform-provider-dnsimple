---
page_title: Migrate from deprecated resources
---

Learn how to migrate resources that have been deprecated or renamed in favour of a new resource.

There can be cases where a resource is deprecated or renamed in a new version of the provider. In these cases, Terraform will not be able to find the resource and will fail to run. In order to migrate to the new resource, you will need to manipulating the state file and configuration files.

~> **Note:** This guide only covers migrating resources that have a replacement resource. If a resource has been removed, you will need to manually remove the resource from the state file.

The steps to migrate a resource are as follows:

1. Remove the resource from the state file
2. Update the configuration files
3. Import the resource into the state file

### Example:

For the purpose of this guide, we will be using `dnsimple_record` as the old resource and `dnsimple_zone_record` as the new resource.

Old resource configuration:

```hcl
locals {
  vegan_pizza = "vegan.pizza"
}

resource "dnsimple_record" "demo" {
  domain = local.vegan_pizza
  name      = "demo"
  value     = "2.3.4.5"
  type      = "A"
  ttl       = 3600
}

...
```

To prepare for the migration we will first want to ensure we have the resource ID. The resource ID can be found in the state or in the DNSimple UI. Refer to each resource's documentation for more information.

To retrieve the ID from the state file you can use the following script that runs code in the `terraform console` and outputs the ID:

```shell
echo dnsimple_record.demo.id | terraform console
```

-> **Note:** To list all the resources tracked in the state, you can use the `terraform state list` command.

Now we need to remove the resource from the state file. To do this, we will need to use the `terraform state rm` command. The command will look something like this:

```shell
terraform state rm dnsimple_record.demo

Removed dnsimple_record.demo
Successfully removed 1 resource instance(s).
```

Update the configuration files to use the new resource:

```hcl
locals {
  vegan_pizza = "vegan.pizza"
}

resource "dnsimple_zone_record" "demo" {
  zone_name = local.vegan_pizza
  name      = "demo"
  value     = "2.3.4.5"
  type      = "A"
  ttl       = 3600
}

...
```

Now we need to import the resource into the state. To do this, we will use the resource ID we previously obtained, and then we will use the [terraform import command](https://www.terraform.io/docs/import/index.html).

To import the resource, we will need to run the following command:

```shell
terraform import dnsimple_zone_record.demo vegan.pizza_2879253

dnsimple_zone_record.demo: Importing from ID "vegan.pizza_2879253"...
dnsimple_zone_record.demo: Import prepared!
  Prepared dnsimple_zone_record for import
dnsimple_zone_record.demo: Refreshing state... [id=2879253]

Import successful!

The resources that were imported are shown above. These resources are now in
your Terraform state and will henceforth be managed by Terraform.
```

-> **Note:** The resource ID for `dnsimple_zone_record` is in the [format](https://registry.terraform.io/providers/dnsimple/dnsimple/latest/docs/resources/zone_record#import) `<zone_name>_<record_id>`. For example, `vegan.pizza_2645561`.

Once the resource has been imported, you can run `terraform plan` and no changes should be reported.

```
terraform plan
dnsimple_record.www: Refreshing state... [id=2879254]
dnsimple_zone_record.demo: Refreshing state... [id=2879253]

No changes. Your infrastructure matches the configuration.

Terraform has compared your real infrastructure against your configuration and found no differences, so no changes are needed.
```
