# confluence custom-content list

List Cloud custom content by type

## Synopsis

List Cloud custom content by type using the documented v2 global endpoint.

## Examples

confluence custom-content list --type ac:example
  confluence custom-content list --type ac:example --space-id 100 --limit 25 --json
  confluence custom-content list --type ac:example --id 777 --body-format storage

## Usage

```text
confluence custom-content list [flags]
```

## Options

```text
      --body-format string   Cloud body representation to include, e.g. storage
      --id strings           Custom content id filter; repeatable or comma-separated
      --limit int            Max results (hard cap 200) (default 25)
      --sort string          Cloud sort expression
      --space-id strings     Space id filter; repeatable or comma-separated
      --type string          Required custom content type
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
