# confluence restriction list

List read/update restrictions on a content id

## Synopsis

List read/update restrictions on a content id.

## Examples

confluence restriction list --page 12345
  confluence restriction list --page 12345 --operation read
  confluence restriction list --page 12345 --operation update --json

## Usage

```text
confluence restriction list [flags]
```

## Options

```text
      --limit int          Maximum users/groups returned for --operation (default 25)
      --operation string   Restriction operation: read or update
      --page string        Content id (required)
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

- [confluence restriction](confluence_restriction.md)
