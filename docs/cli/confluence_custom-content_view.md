# confluence custom-content view

Show one Cloud custom-content record

## Synopsis

Show one Cloud custom-content record by id.

## Examples

confluence custom-content view 12345
  confluence custom-content view 12345 --body-format storage --include-version --json
  confluence custom-content view 12345 --include-labels --include-properties

## Usage

```text
confluence custom-content view <id> [flags]
```

## Options

```text
      --body-format string      Cloud body representation to include, e.g. storage
      --include-collaborators   Include collaborators in the Cloud response
      --include-labels          Include labels in the Cloud response
      --include-operations      Include permitted operations in the Cloud response
      --include-properties      Include properties in the Cloud response
      --include-version         Include current version in the Cloud response
      --include-versions        Include versions in the Cloud response
      --version int             Cloud content version to retrieve
```

## Inherited Options

```text
      --debug             Log HTTP requests to stderr (Authorization header redacted)
      --jq string         Filter JSON output through a jq expression (requires jq on PATH)
      --json              Output as JSON (machine-readable)
      --template string   Render JSON output through a Go text/template
      --timing            Show execution time on stderr
```

## See Also

- [confluence custom-content](confluence_custom-content.md)
