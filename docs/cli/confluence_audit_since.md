# confluence audit since

List recent Cloud audit records

## Synopsis

List Cloud audit records for a time period back from the current date.

## Examples

confluence audit since --number 3 --unit MONTHS
  confluence audit since --number 7 --unit DAYS --search group --json

## Usage

```text
confluence audit since [flags]
```

## Options

```text
      --limit int       Max records (hard cap 200) (default 25)
      --number int      Cloud time period number
      --search string   Cloud audit search string
      --unit string     Cloud time period unit, e.g. DAYS, WEEKS, MONTHS (default "MONTHS")
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

- [confluence audit](confluence_audit.md)
