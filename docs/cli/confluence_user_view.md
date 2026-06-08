# confluence user view

Show a user by username / key / accountId

## Synopsis

Show a user. Pick one selector. On Server/DC use --username or --key;
on Cloud use --accountId.

## Examples

confluence user view --username alice
  confluence user view --accountId 557058:abc
  confluence user view --key u1234

## Usage

```text
confluence user view [flags]
```

## Options

```text
      --accountId string   Cloud accountId
      --key string         Server/DC user key
      --username string    Server/DC username
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
