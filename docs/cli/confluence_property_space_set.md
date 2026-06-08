# confluence property space set

Create or update one space property

## Synopsis

Create or update a space property. --value must be JSON.
Use JSON strings for scalar text values, for example --value '"platform"'.

## Examples

confluence property space set --space ENG --key retention --value '{"days":30}'
  confluence property space set --space ENG --key owner --value '"platform"'
  confluence property space set --space ENG --key retention --value @property.json

## Usage

```text
confluence property space set [flags]
```

## Options

```text
      --key string     Property key (required)
      --space string   Space key (required)
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

- [confluence property space](confluence_property_space.md)
