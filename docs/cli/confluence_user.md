# confluence user

Users (current, view, search, bulk)

## Synopsis

User operations.

## Examples

confluence user current
  confluence user view --username alice                  # Server/DC
  confluence user view --account-id 557058:abc...         # Cloud
  confluence user search "Jane Smith"
  confluence user bulk --account-id 557058:abc --json

## Usage

```text
confluence user
```

## Commands

- `bulk` - Show Cloud users by account id
- `current` - Show the authenticated user
- `search` - Search users by full name (CQL user.fullname~)
- `view` - Show a user by username / key / account id

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
