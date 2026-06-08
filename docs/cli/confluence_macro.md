# confluence macro

Macro body utilities

## Synopsis

Macro body utility operations.

## Examples

confluence macro body --page 12345 --version 2 --macro-id 50884bd9-0cb8-41d5-98be-f80943c14f96
  confluence macro body --page 12345 --version 2 --hash abc123 --json

## Usage

```text
confluence macro
```

## Commands

- `body` - Fetch one macro body

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
