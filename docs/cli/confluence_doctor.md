# confluence doctor

Verify environment, auth, and flavor detection

## Verify that

1. Required environment variables are set
  2. Flavor is detected (server | cloud)
  3. The auth credentials are valid (GET /rest/api/user/current)

Exits 0 on success, 2 on user-fixable config errors, 1 on unexpected errors.

## Examples

confluence doctor
  confluence doctor --json

## Usage

```text
confluence doctor
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

- [confluence](confluence.md)
