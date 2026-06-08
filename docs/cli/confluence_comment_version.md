# confluence comment version

Show one comment version

## Synopsis

Show one version record for one comment.

Cloud uses the documented v2 Version detail endpoints. Server/Data Center uses
the documented content version detail route where the target is content-like.

## Examples

confluence comment version 12345 2
  confluence comment version 12345 2 --json

## Usage

```text
confluence comment version <id> <number> [flags]
```

## Options

```text
      --location string   Cloud comment location: footer or inline (default "footer")
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

- [confluence comment](confluence_comment.md)
