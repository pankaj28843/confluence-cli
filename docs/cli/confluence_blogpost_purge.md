# confluence blogpost purge

Permanently delete a trashed blog post

## Synopsis

Permanently delete a trashed blog post. Cloud requires the blog post to
already be in trash. Server/Data Center sends status=trashed to purge trashable
content.

## Examples

confluence blogpost purge 12345
  confluence blogpost purge 12345 --force
  confluence blogpost purge 12345 --force --json

## Usage

```text
confluence blogpost purge <id> [flags]
```

## Options

```text
      --force   Do not prompt for confirmation
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
