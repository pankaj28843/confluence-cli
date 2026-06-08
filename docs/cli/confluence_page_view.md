# confluence page view

Fetch a page (or other content) by id

## Synopsis

Fetch content by id. --markdown renders the storage body to Markdown;
--raw-storage emits the raw Confluence storage XML/HTML. --body-only emits only
raw storage HTML for edit-and-reupload workflows.

## Examples

confluence page view 12345 --markdown
  confluence page view 12345 --json
  confluence page view 12345 --raw-storage
  confluence page view 12345 --body-only > body.html

## Usage

```text
confluence page view <id> [flags]
```

## Options

```text
      --body-only       Emit only raw storage-format XML/HTML
      --expand string   Override the expand= parameter (default: body.storage,version,space,ancestors,metadata.labels)
      --markdown        Render body as Markdown (default output)
      --raw-storage     Emit the raw storage-format XML
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
