---
page_title: Migrate from deprecated resources
---

# Migrate from Deprecated Resources

Learn how to migrate resources that have been deprecated or renamed in favor of a new resource.

## Overview

When a resource is deprecated or renamed in a new version of the provider, Terraform will not be able to find the resource and will fail to run. To migrate to the new resource, you will need to manipulate the state file and configuration files.

-> **Note:** This guide only covers migrating resources that have a replacement resource. If a resource has been removed without a replacement, you will need to manually remove the resource from the state file using `terraform state rm`.

## Migration Steps

The migration process consists of three main steps:

1. **Obtain the resource ID** - Identify the unique identifier for the resource you're migrating
2. **Remove the old resource from state** - Remove the deprecated resource from Terraform's state file
3. **Update configuration and import** - Update your configuration files to use the new resource and import it into state

## Example Migration

For the purpose of this guide, we will demonstrate migrating from `dnsimple_record` (deprecated) to `dnsimple_zone_record` (current).

### Step 1: Obtain the Resource ID

Before starting the migration, you need to identify the resource ID. The resource ID can be found in several ways:

- **From the Terraform state file** - Use `terraform console` to query the state
- **From the DNSimple UI** - View the resource in the DNSimple web dashboard
- **From the DNSimple API** - Query the API directly

To retrieve the ID from the state file, you can use the following command:

```shell
echo dnsimple_record.demo.id | terraform console
```

-> **Note:** To list all the resources tracked in the state, you can use the `terraform state list` command.

For `dnsimple_zone_record`, the import ID format is `<zone_name>_<record_id>`. You may need to construct this ID from separate values. Refer to each resource's documentation for the specific import ID format.

### Step 2: Remove the Old Resource from State

Remove the deprecated resource from Terraform's state file using the `terraform state rm` command:

```shell
terraform state rm dnsimple_record.demo
```

Expected output:

```
Removed dnsimple_record.demo
Successfully removed 1 resource instance(s).
```

-> **Important:** This command only removes the resource from Terraform's state file. It does not delete the actual resource in DNSimple. The resource will continue to exist and function normally.

### Step 3: Update Configuration and Import

Update your configuration files to use the new resource. Note any differences in attribute names or required fields between the old and new resources.

**Old resource configuration:**

```hcl
locals {
  vegan_pizza = "vegan.pizza"
}

resource "dnsimple_record" "demo" {
  domain = local.vegan_pizza
  name   = "demo"
  value  = "2.3.4.5"
  type   = "A"
  ttl    = 3600
}
```

**New resource configuration:**

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
```

Notice that `domain` has been changed to `zone_name` in the new resource.

Now import the resource into the state using the `terraform import` command with the resource ID you obtained earlier:

```shell
terraform import dnsimple_zone_record.demo vegan.pizza_2879253
```

Expected output:

```
dnsimple_zone_record.demo: Importing from ID "vegan.pizza_2879253"...
dnsimple_zone_record.demo: Import prepared!
  Prepared dnsimple_zone_record for import
dnsimple_zone_record.demo: Refreshing state... [id=2879253]

Import successful!

The resources that were imported are shown above. These resources are now in
your Terraform state and will henceforth be managed by Terraform.
```

-> **Note:** The resource ID for `dnsimple_zone_record` is in the format `<zone_name>_<record_id>`. For example, `vegan.pizza_2645561`. Refer to the [resource documentation](https://registry.terraform.io/providers/dnsimple/dnsimple/latest/docs/resources/zone_record#import) for more details.

### Step 4: Verify the Migration

After importing the resource, run `terraform plan` to verify that Terraform recognizes the resource and that no changes are needed:

```shell
terraform plan
```

Expected output:

```
dnsimple_zone_record.demo: Refreshing state... [id=2879253]

No changes. Your infrastructure matches the configuration.

Terraform has compared your real infrastructure against your configuration and found no differences, so no changes are needed.
```

If `terraform plan` shows differences, review the configuration to ensure all attributes match the imported resource. You may need to adjust attribute values or add missing attributes.

## Troubleshooting

### Import ID Format

If the import fails, verify that you're using the correct import ID format for the resource. Each resource type has a specific format documented in its resource page. Common formats include:

- `<zone_name>_<record_id>` for zone records
- Domain names for domain resources
- Numeric IDs for other resources

### Configuration Mismatches

If `terraform plan` shows changes after import, compare the imported resource attributes with your configuration. Common issues include:

- Attribute name changes (e.g., `domain` → `zone_name`)
- Default value differences
- Missing optional attributes that were set in the original resource

### Multiple Resources

When migrating multiple resources, repeat the process for each resource. You can automate this with scripts, but be careful to verify each migration individually.

## Additional Resources

- [Terraform Import Documentation](https://www.terraform.io/docs/import/index.html) - Official Terraform import guide
- [Terraform State Management](https://www.terraform.io/docs/state/index.html) - Understanding Terraform state
- [DNSimple Provider Documentation](https://registry.terraform.io/providers/dnsimple/dnsimple/latest/docs) - Complete provider reference
