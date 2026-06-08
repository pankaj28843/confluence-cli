# confluence watcher status

Show whether a user watches content or a space

## Synopsis

Show whether the current or specified user watches a content id or
space key.

## Examples

confluence watcher status --page 12345
  confluence watcher status --space ENG --content-type blogpost --json
  confluence watcher status --page 12345 --account-id abc123

## Usage

```text
confluence watcher status [flags]
```

## Options

```text
      --account-id string     Cloud account id to check
      --content-type string   Space watch content type, such as page or blogpost
      --page string           Content id
      --space string          Space key
      --user-key string       Server/DC user key to check
      --username string       Server/DC username to check
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

- [confluence watcher](confluence_watcher.md)
