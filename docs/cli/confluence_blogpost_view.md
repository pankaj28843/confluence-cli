# confluence blogpost view

Fetch a blog post by id

## Synopsis

Fetch a blog post by id. --markdown renders the storage body to Markdown;
--raw-storage emits raw Confluence storage XML/HTML. --body-only emits only the
storage body for edit-and-reupload workflows.

## Examples

confluence blogpost view 12345 --markdown
  confluence blogpost view 12345 --json
  confluence blogpost view 12345 --body-only > body.html

## Usage

```text
confluence blogpost view <id> [flags]
```

## Options

```text
      --body-only     Emit only raw storage-format XML/HTML
      --markdown      Render body as Markdown (default output)
      --raw-storage   Emit the raw storage-format XML
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
