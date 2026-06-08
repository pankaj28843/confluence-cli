# confluence custom-content page

List Cloud custom content in a page

## Synopsis

List Cloud custom content by type inside one page.

## Examples

confluence custom-content page 12345 --type ac:example
  confluence custom-content page 12345 --type ac:example --body-format storage --json

## Usage

```text
confluence custom-content page <id> [flags]
```

## Options

```text
      --body-format string   Cloud body representation to include, e.g. storage
      --limit int            Max results (hard cap 200) (default 25)
      --sort string          Cloud sort expression
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
