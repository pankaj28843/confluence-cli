# confluence blogpost version

Show one blogpost version

## Synopsis

Show one version record for one blogpost.

Cloud uses the documented v2 Version detail endpoints. Server/Data Center uses
the documented content version detail route where the target is content-like.

## Examples

confluence blogpost version 12345 2
  confluence blogpost version 12345 2 --json

## Usage

```text
confluence blogpost version <id> <number>
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

- [confluence blogpost](confluence_blogpost.md)
