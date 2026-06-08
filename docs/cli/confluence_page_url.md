# confluence page url

Print the browser URL for a page

## Synopsis

Print the absolute browser URL for a page using Confluence's _links.base
and _links.webui fields.

## Examples

confluence page url 12345
  cdp open "$(confluence page url 12345)" --new-tab=false

## Usage

```text
confluence page url <id>
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

- [confluence page](confluence_page.md)
