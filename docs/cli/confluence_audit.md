# confluence audit

Audit log reads

## Synopsis

Audit log read operations.

## Examples

confluence audit list --limit 25 --json
  confluence audit since --number 7 --unit DAYS --search group
  confluence audit retention --json

## Usage

```text
confluence audit
```

## Commands

- `list` - List audit records
- `retention` - Show Cloud audit retention period
- `since` - List recent Cloud audit records

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
