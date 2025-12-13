---
page_title: Import existing resources
---

# Import Existing Resources

Learn how to import existing DNSimple resources into Terraform so you can manage them as infrastructure as code.

## Overview

When you have existing resources in DNSimple that were created outside of Terraform, you can import them into your Terraform state. This allows Terraform to manage these resources going forward without recreating them.

Terraform provides two methods for importing resources:

1. **`import` block** - Configuration-driven approach (recommended for remote state backends like Terraform Cloud)
2. **`terraform import` command** - Traditional command-line approach

-> **Note:** The `import` block method is recommended when using remote state backends like Terraform Cloud, as it allows the import process to be managed entirely through code without requiring manual command execution.

## Import Process

The import process consists of three main steps:

1. **Extract resource information** - Retrieve the resource details from DNSimple using the API
2. **Define the resource in Terraform** - Create the resource configuration matching the existing resource
3. **Import into state** - Import the resource using either the `terraform import` command or an `import` block

## Example: Importing a Contact

This guide demonstrates importing a DNSimple contact resource. The same process applies to other resource types, though the API endpoints and import ID formats will differ.

### Step 1: Extract Resource Information from DNSimple

Before importing, you need to retrieve the resource details from DNSimple. This helps you understand the current configuration and ensures your Terraform configuration matches the existing resource.

#### List All Contacts

To see all contacts in your account, use the DNSimple API:

```shell
curl -s -H "Authorization: Bearer $DNSIMPLE_TOKEN" \
     -H "Accept: application/json" \
     "https://api.dnsimple.com/v2/$DNSIMPLE_ACCOUNT/contacts" | jq
```

This will return a JSON response listing all contacts:

```json
{
  "data": [
    {
      "id": 12345,
      "account_id": 1010,
      "label": "Main Contact",
      "first_name": "John",
      "last_name": "Doe",
      "organization_name": "Example Inc",
      "job_title": "Manager",
      "address1": "123 Main Street",
      "address2": "Suite 100",
      "city": "San Francisco",
      "state_province": "California",
      "postal_code": "94105",
      "country": "US",
      "phone": "+1.4155551234",
      "fax": "+1.4155555678",
      "email": "john@example.com",
      "phone_normalized": "+1.4155551234",
      "fax_normalized": "+1.4155555678",
      "created_at": "2023-01-15T10:30:00Z",
      "updated_at": "2023-01-15T10:30:00Z"
    }
  ]
}
```

#### Get a Specific Contact

If you know the contact ID, you can retrieve a specific contact's details:

```shell
curl -s -H "Authorization: Bearer $DNSIMPLE_TOKEN" \
     -H "Accept: application/json" \
     "https://api.dnsimple.com/v2/$DNSIMPLE_ACCOUNT/contacts/12345" | jq
```

This returns the contact details in a similar format:

```json
{
  "data": {
    "id": 12345,
    "account_id": 1010,
    "label": "Main Contact",
    "first_name": "John",
    "last_name": "Doe",
    "organization_name": "Example Inc",
    "job_title": "Manager",
    "address1": "123 Main Street",
    "address2": "Suite 100",
    "city": "San Francisco",
    "state_province": "California",
    "postal_code": "94105",
    "country": "US",
    "phone": "+1.4155551234",
    "fax": "+1.4155555678",
    "email": "john@example.com",
    "phone_normalized": "+1.4155551234",
    "fax_normalized": "+1.4155555678",
    "created_at": "2023-01-15T10:30:00Z",
    "updated_at": "2023-01-15T10:30:00Z"
  }
}
```

-> **Note:** Replace `$DNSIMPLE_TOKEN` and `$DNSIMPLE_ACCOUNT` with your actual API token and account ID, or set them as environment variables. The contact ID (`12345` in this example) is the value you'll use for the import.

### Step 2: Define the Resource in Terraform

Create or update your Terraform configuration file to define the contact resource. Use the information retrieved from the API to match the existing resource configuration.

```hcl
resource "dnsimple_contact" "main" {
  label            = "Main Contact"
  first_name       = "John"
  last_name        = "Doe"
  organization_name = "Example Inc"
  job_title        = "Manager"
  address1         = "123 Main Street"
  address2         = "Suite 100"
  city             = "San Francisco"
  state_province   = "California"
  postal_code      = "94105"
  country          = "US"
  phone            = "+1.4155551234"
  fax              = "+1.4155555678"
  email            = "john@example.com"
}
```

-> **Important:** Ensure your Terraform configuration matches the existing resource attributes. After importing, run `terraform plan` to verify there are no differences. If there are differences, update your configuration accordingly.

### Step 3: Import the Resource

You can import the resource using either an `import` block or the `terraform import` command. The `import` block method is recommended for remote state backends.

#### Method 1: Using `import` Block (Recommended)

The `import` block allows you to define imports directly in your Terraform configuration. This is particularly useful when using remote state backends like Terraform Cloud, as it eliminates the need to run manual import commands.

Add an `import` block to your Terraform configuration:

```hcl
import {
  to = dnsimple_contact.main
  id = "12345"
}

resource "dnsimple_contact" "main" {
  label            = "Main Contact"
  first_name       = "John"
  last_name        = "Doe"
  organization_name = "Example Inc"
  job_title        = "Manager"
  address1         = "123 Main Street"
  address2         = "Suite 100"
  city             = "San Francisco"
  state_province   = "California"
  postal_code      = "94105"
  country          = "US"
  phone            = "+1.4155551234"
  fax              = "+1.4155555678"
  email            = "john@example.com"
}
```

Then run `terraform plan`:

```shell
terraform plan
```

Terraform will detect the `import` block and show the import operation:

```
dnsimple_contact.main: Preparing import... [id=12345]
dnsimple_contact.main: Refreshing state... [id=12345]

Terraform will perform the following actions:

  # dnsimple_contact.main will be imported
  <= resource "dnsimple_contact" "main" {
      ...
    }

Plan: 0 to add, 0 to change, 0 to destroy, 1 to import.
```

Apply the import:

```shell
terraform apply
```

After the import is complete, you can remove the `import` block from your configuration. The resource is now managed by Terraform and the import block is no longer needed.

-> **Note:** The `import` block is available in Terraform 1.5 and later. This method is recommended when using remote state backends because it allows the import to be managed through your configuration and version control, making it easier to track and reproduce.

#### Method 2: Using `terraform import` Command

Use the `terraform import` command to import the contact into your Terraform state:

```shell
terraform import dnsimple_contact.main 12345
```

Expected output:

```
dnsimple_contact.main: Importing from ID "12345"...
dnsimple_contact.main: Import prepared!
  Prepared dnsimple_contact for import
dnsimple_contact.main: Refreshing state... [id=12345]

Import successful!

The resources that were imported are shown above. These resources are now in
your Terraform state and will henceforth be managed by Terraform.
```

After importing, run `terraform plan` to verify the configuration matches the imported resource:

```shell
terraform plan
```

If the configuration matches, you should see:

```
No changes. Your infrastructure matches the configuration.
```

## Import ID Formats

Different resource types use different import ID formats. Refer to each resource's documentation for the specific format:

- **Contacts** - Numeric ID only (e.g., `12345`)
- **Zone Records** - Format: `<zone_name>_<record_id>` (e.g., `example.com_2645561`)
- **Domains** - Domain name (e.g., `example.com`)
- **Registered Domains** - Domain name (e.g., `example.com`)

Refer to the individual resource documentation pages for complete import instructions and ID formats.

## Troubleshooting

### Import Fails with "Resource Not Found"

If the import fails with a "resource not found" error:

- Verify the resource ID is correct
- Ensure you're using the correct account ID in your provider configuration
- Check that the resource exists in DNSimple using the API

### Configuration Mismatches After Import

If `terraform plan` shows differences after import:

- Compare the imported resource attributes with your configuration
- Check for optional attributes that may have default values
- Verify attribute names match exactly (some attributes may have changed between provider versions)
- Review computed attributes that may differ (e.g., normalized phone numbers)

### Multiple Resources

When importing multiple resources:

- Import one resource at a time to verify each import is successful
- Use `terraform plan` after each import to check for configuration mismatches
- Consider using `import` blocks for multiple resources to manage them all in code

## Additional Resources

- [Terraform Import Documentation](https://developer.hashicorp.com/terraform/cli/import) - Official Terraform import guide
- [Terraform Import Blocks](https://developer.hashicorp.com/terraform/language/import) - Documentation on using import blocks
- [DNSimple API Documentation](https://developer.dnsimple.com/) - Complete API reference
- [DNSimple Provider Documentation](https://registry.terraform.io/providers/dnsimple/dnsimple/latest/docs) - Complete provider reference

