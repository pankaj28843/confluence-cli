# confluence permission space available

List available Cloud space permissions

## Synopsis

List available Cloud space permissions.

This is the documented Cloud v2 permission catalog for tenants with RBAC.

## Examples

confluence permission space available --json
  confluence permission space available --limit 100

## Usage

```text
confluence permission space available [flags]
```

## Options

```text
      --limit int   Max permissions (hard cap 200) (default 25)
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

- [confluence permission space](confluence_permission_space.md)
