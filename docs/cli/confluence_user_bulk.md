# confluence user bulk

Show Cloud users by account id

## Synopsis

Show Cloud users by account id using the documented v2 users-bulk endpoint.

## Examples

confluence user bulk --account-id 557058:abc
  confluence user bulk --account-id 557058:abc --account-id 557058:def --json

## Usage

```text
confluence user bulk [flags]
```

## Options

```text
      --account-id strings   Cloud account id; repeatable or comma-separated
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

- [confluence user](confluence_user.md)
