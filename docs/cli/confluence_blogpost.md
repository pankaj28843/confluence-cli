# confluence blogpost

Blog posts (list, view, create, update, delete, purge)

## Synopsis

Blog post operations.

## Examples

confluence blogpost list --space ENG
  confluence blogpost view 12345 --markdown
  confluence blogpost create --space ENG --title "Weekly Update" --body-file body.html

## Usage

```text
confluence blogpost
```

## Commands

- `create` - Create a new blog post
- `delete` - Move a blog post to trash
- `list` - List blog posts
- `purge` - Permanently delete a trashed blog post
- `update` - Update an existing blog post
- `view` - Fetch a blog post by id

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
