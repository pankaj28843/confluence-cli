# confluence comment update

Update a footer comment body

## Synopsis

Update a footer comment body. The command fetches the current comment version
and writes version.number + 1 unless --version is supplied.

## Examples

confluence comment update 998877 --body "<p>Updated.</p>"
  confluence comment update 998877 --body-file comment.html
  echo "<p>Updated.</p>" | confluence comment update 998877 --body-file - --json

## Usage

```text
confluence comment update <id> [flags]
```

## Options

```text
      --body string          Inline body string
      --body-file string     Path to body file, or '-' for stdin
      --body-format string   Body format: storage (default "storage")
      --version string       Explicit current version number (auto-fetched if omitted)
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

- [confluence comment](confluence_comment.md)
