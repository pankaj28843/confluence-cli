# confluence audit list

List audit records

## Synopsis

List audit records. Cloud supports date and search filters; Server/Data Center uses the documented deprecated read endpoint.

## Examples

confluence audit list
  confluence audit list --start-date 1700000000000 --end-date 1700100000000 --search space --json

## Usage

```text
confluence audit list [flags]
```

## Options

```text
      --end-date string     Cloud end date as epoch milliseconds
      --limit int           Max records (hard cap 200) (default 25)
      --search string       Cloud audit search string
      --start-date string   Cloud start date as epoch milliseconds
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
