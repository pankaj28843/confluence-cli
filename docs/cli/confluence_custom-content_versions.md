# confluence custom-content versions

List custom-content versions

## Synopsis

List version records for one custom-content.

Cloud uses the documented v2 Version endpoints. Server/Data Center uses the
documented content version route where the target is content-like.

## Examples

confluence custom-content versions 12345
  confluence custom-content versions 12345 --limit 25 --json
  confluence custom-content versions 12345 --body-format storage --sort -modified-date

## Usage

```text
confluence custom-content versions <id> [flags]
```

## Options

```text
      --body-format string   Cloud body representation to include, e.g. storage
      --limit int            Max versions (hard cap 200) (default 25)
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

- [confluence custom-content](confluence_custom-content.md)
