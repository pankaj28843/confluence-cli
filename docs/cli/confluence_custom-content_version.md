# confluence custom-content version

Show one custom-content version

## Synopsis

Show one version record for one custom-content.

Cloud uses the documented v2 Version detail endpoints. Server/Data Center uses
the documented content version detail route where the target is content-like.

## Examples

confluence custom-content version 12345 2
  confluence custom-content version 12345 2 --json

## Usage

```text
confluence custom-content version <id> <number>
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
