# confluence blogpost create

Create a new blog post

## Synopsis

Create a new blog post with a storage-format body. On Cloud this uses
the v2 blogposts endpoint; on Server/Data Center this uses the content resource
with type=blogpost.

## Examples

confluence blogpost create --space ENG --title "Weekly Update" --body-file body.html
  confluence blogpost create --space ENG --title "Draft" --draft --body "<p>Hello</p>"
  echo "<p>Hello</p>" | confluence blogpost create --space ENG --title "Hello" --body-file -

## Usage

```text
confluence blogpost create [flags]
```

## Options

```text
      --body string          Inline body string
      --body-file string     Path to body file, or '-' for stdin
      --body-format string   Body format: storage | wiki | view (default "storage")
      --draft                Create as draft
      --space string         Space key (required)
      --title string         Blog post title (required)
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
