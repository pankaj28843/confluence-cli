# confluence watcher

Watchers (content, space)

## Synopsis

Watcher operations.

## Examples

confluence watcher content --page 12345
  confluence watcher space --space ENG

## Usage

```text
confluence watcher
```

## Commands

- `content` - Watchers subscribed to a content id (raw passthrough)
- `space` - Watchers subscribed to a space (raw passthrough)

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
