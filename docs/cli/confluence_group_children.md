# confluence group children

List Server/DC child groups

## Synopsis

List Server/DC child groups for a Server/Data Center group.

## Examples

confluence group children engineering
  confluence group children engineering --expand members --limit 50 --json

## Usage

```text
confluence group children <name> [flags]
```

## Options

```text
      --expand string   Server/DC expand value
      --limit int       Max groups (hard cap 200) (default 25)
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

- [confluence group](confluence_group.md)
