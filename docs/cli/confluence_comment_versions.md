# confluence comment versions

List comment versions

## Synopsis

List version records for one comment.

Cloud uses the documented v2 Version endpoints. Server/Data Center uses the
documented content version route where the target is content-like.

## Examples

confluence comment versions 12345
  confluence comment versions 12345 --limit 25 --json
  confluence comment versions 12345 --body-format storage --sort -modified-date

## Usage

```text
confluence comment versions <id> [flags]
```

## Options

```text
      --body-format string   Cloud body representation to include, e.g. storage
      --limit int            Max versions (hard cap 200) (default 25)
      --location string      Cloud comment location: footer or inline (default "footer")
      --sort string          Cloud sort expression
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
