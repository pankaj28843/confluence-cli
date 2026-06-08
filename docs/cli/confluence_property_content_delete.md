# confluence property content delete

Delete one content property by key

## Synopsis

Delete a property for a page/content id by key.

## Examples

confluence property content delete --page 12345 --key release
  confluence property content delete --page 12345 --key release --json

## Usage

```text
confluence property content delete [flags]
```

## Options

```text
      --key string    Property key (required)
      --page string   Content id (required)
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

- [confluence property content](confluence_property_content.md)
