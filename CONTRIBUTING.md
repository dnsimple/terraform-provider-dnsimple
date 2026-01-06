# Contributing to DNSimple/terraform-provider

## Contributing Workflow

1. **Create a new branch**

   Create a feature branch from `main` for your changes:

   ```shell
   git checkout -b feature/your-feature-name
   ```

2. **Make the changes**

   Implement your changes, following the code style guidelines below.

3. **Validate tests and linters**

   Before submitting your PR, ensure all checks pass:

   ```shell
   make fmtcheck    # Check code formatting
   make errcheck    # Check for unchecked errors
   make test        # Run unit tests
   ```

   If formatting issues are found, you can auto-fix them:

   ```shell
   make fmt         # Format code with gofumpt
   ```

   If you've modified any generated code, ensure code generation is up to date:

   ```shell
   go generate ./...
   ```

4. **Create the PR**

   Push your branch and create a pull request against `main`. Include a clear description of your changes and reference any related issues.

5. **Follow up**

   - Respond to any review feedback promptly
   - Make requested changes and push updates to your branch
   - Ensure CI checks pass (tests, formatting, and static analysis)

## Changelog

We loosely follow the [Common Changelog](https://common-changelog.org/) format for changelog entries.

## Code Style and Static Analysis

We use several tools to maintain code quality and consistency:

### Code Formatting

We use [`gofumpt`](https://github.com/mvdan/gofumpt) for code formatting, which is a stricter version of `gofmt`.

- **Check formatting**: Run `make fmtcheck` to verify your code is properly formatted
- **Auto-format**: Run `make fmt` to automatically format your code

The build and test targets automatically run `fmtcheck` to ensure all code is properly formatted.

### Error Checking

We use [`errcheck`](https://github.com/kisielk/errcheck) to ensure all errors are properly handled.

- **Check for unchecked errors**: Run `make errcheck`

This helps prevent bugs by ensuring all function return values, especially errors, are properly handled.

## Testing

Submit unit tests for your changes. You can test your changes on your machine by [running the test suite](README.md#testing):

```shell
make test
```

When you submit a PR, tests will also be run on the continuous integration environment [via GitHub Actions](https://github.com/dnsimple/terraform-provider-dnsimple/actions).
