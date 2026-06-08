# confluence property content list

List content properties

## Synopsis

List properties for a page/content id.

## Examples

confluence property content list --page 12345
  confluence property content list --page 12345 --key release --json
  confluence property content list --page 12345 --limit 100

## Usage

```text
confluence property content list [flags]
```

## Options

```text
      --key string    Optional property key filter
      --limit int     Max properties (hard cap 200) (default 25)
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
