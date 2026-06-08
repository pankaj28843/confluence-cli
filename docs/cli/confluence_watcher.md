# confluence watcher

Watchers (content, space, status)

## Synopsis

Watcher read helpers.

## Examples

confluence watcher content --page 12345
  confluence watcher space --space ENG
  confluence watcher status --page 12345 --json

## Usage

```text
confluence watcher
```

## Commands

- `content` - List watchers subscribed to a content id
- `space` - List watchers subscribed to a space
- `status` - Show whether a user watches content or a space

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
