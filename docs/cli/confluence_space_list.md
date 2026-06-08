# confluence space list

List spaces in the site

## Synopsis

List spaces with optional filters.

## Examples

confluence space list --json --limit 100
  confluence space list --type global --status current

## Usage

```text
confluence space list [flags]
```

## Options

```text
      --limit int       Max spaces to return (hard cap 200) (default 25)
      --status string   Filter: current | archived
      --type string     Filter: global | personal
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

- [confluence space](confluence_space.md)
