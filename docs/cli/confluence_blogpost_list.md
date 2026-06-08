# confluence blogpost list

List blog posts

## Synopsis

List blog posts. On Cloud, --space resolves a space key to its v2 space id.
On Server/Data Center, --posting-day maps to the documented content resource
postingDay filter.

## Examples

confluence blogpost list --space ENG
  confluence blogpost list --space ENG --title "Weekly" --limit 10 --json
  confluence blogpost list --posting-day 2026-06-08

## Usage

```text
confluence blogpost list [flags]
```

## Options

```text
      --label-id string      Cloud only: label id filter
      --limit int            Max results (hard cap 200) (default 25)
      --posting-day string   Server/Data Center only: posting day YYYY-MM-DD
      --space string         Space key filter
      --status string        Status filter, e.g. current,draft,trashed
      --title string         Exact title filter
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
