# confluence property

Content and space properties

## Synopsis

Content and space property operations.

## Examples

confluence property content list --page 12345
  confluence property content get --page 12345 --key release
  confluence property space set --space ENG --key retention --value '{"days":30}'

## Usage

```text
confluence property
```

## Commands

- `content` - Content properties (list, get, set, delete)
- `space` - Space properties (list, get, set, delete)

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
