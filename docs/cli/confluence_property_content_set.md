# confluence property content set

Create or update one content property

## Synopsis

Create or update a property for a page/content id. --value must be JSON.
Use JSON strings for scalar text values, for example --value '"done"'.

## Examples

confluence property content set --page 12345 --key release --value '{"ready":true}'
  confluence property content set --page 12345 --key owner --value '"platform"'
  confluence property content set --page 12345 --key release --value @property.json

## Usage

```text
confluence property content set [flags]
```

## Options

```text
      --key string     Property key (required)
      --page string    Content id (required)
      --value string   JSON value, @file, or @- (required)
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
