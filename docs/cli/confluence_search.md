# confluence search

Search content, spaces, users, attachments, or all

## Synopsis

Search Confluence. Five sub-verbs — content, spaces, users, attachments, all.
'all' fans out to the first four in parallel and merges via reciprocal-rank fusion.

## Examples

confluence search content "release"
  confluence search spaces "engineering"
  confluence search users "Jane Smith"
  confluence search attachments "logo"
  confluence search all "release process" --json

## Usage

```text
confluence search
```

## Commands

- `all` - Unified search — pages + spaces + users + attachments in parallel
- `attachments` - Search attachments via CQL type=attachment
- `content` - Search pages (type=page by default) via CQL text match
- `spaces` - Search spaces via CQL type=space
- `users` - Search users via CQL user.fullname~

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
