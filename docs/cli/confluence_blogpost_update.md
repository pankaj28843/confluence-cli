# confluence blogpost update

Update an existing blog post

## Synopsis

Update a blog post title, body, or both. The command fetches the current
blog post version and body, then writes version.number + 1 unless --version is
supplied.

## Examples

confluence blogpost update 12345 --title "New Title"
  confluence blogpost update 12345 --body-file body.html
  echo "<p>Hello</p>" | confluence blogpost update 12345 --body-file -

## Usage

```text
confluence blogpost update <id> [flags]
```

## Options

```text
      --body string          Inline body string
      --body-file string     Path to body file, or '-' for stdin
      --body-format string   Body format: storage | wiki | view (default "storage")
      --title string         New title (keeps existing if omitted)
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

- [confluence blogpost](confluence_blogpost.md)
