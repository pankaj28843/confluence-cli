# confluence attachment versions

List attachment versions

## Synopsis

List version records for one attachment.

Cloud uses the documented v2 Version endpoints. Server/Data Center uses the
documented content version route where the target is content-like.

## Examples

confluence attachment versions 12345
  confluence attachment versions 12345 --limit 25 --json
  confluence attachment versions 12345 --sort -modified-date

## Usage

```text
confluence attachment versions <id> [flags]
```

## Options

```text
      --limit int     Max versions (hard cap 200) (default 25)
      --sort string   Cloud sort expression
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

- [confluence attachment](confluence_attachment.md)
