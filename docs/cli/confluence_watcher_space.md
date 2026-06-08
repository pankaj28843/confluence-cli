# confluence watcher space

Watchers subscribed to a space (raw passthrough)

## Synopsis

Show watcher records for a space key.

## Examples

confluence watcher space --space ENG --json

## Usage

```text
confluence watcher space [flags]
```

## Options

```text
      --space string   Space key (required)
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
