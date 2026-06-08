# confluence page purge

Permanently delete a trashed page

## Synopsis

Permanently delete a trashed page. Cloud requires the page to already be in
trash. Server/Data Center sends status=trashed to purge trashable content.

## Examples

confluence page purge 12345
  confluence page purge 12345 --force
  confluence page purge 12345 --force --json

## Usage

```text
confluence page purge <id> [flags]
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

- [confluence page](confluence_page.md)
