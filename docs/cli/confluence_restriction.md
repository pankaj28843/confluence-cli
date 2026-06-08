# confluence restriction

Content restrictions (list)

## Synopsis

Content restriction helpers.

List restrictions grouped by operation, or inspect the read/update restrictions
for one content id.

## Examples

confluence restriction list --page 12345
  confluence restriction list --page 12345 --operation read --json

## Usage

```text
confluence restriction
```

## Commands

- `list` - List read/update restrictions on a content id

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
