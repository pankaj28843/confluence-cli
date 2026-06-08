# confluence property space delete

Delete one space property by key

## Synopsis

Delete a space property by key.

## Examples

confluence property space delete --space ENG --key retention
  confluence property space delete --space ENG --key retention --json

## Usage

```text
confluence property space delete [flags]
```

## Options

```text
      --key string     Property key (required)
      --space string   Space key (required)
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
