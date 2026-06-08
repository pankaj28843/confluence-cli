# confluence page create

Create a new page

## Synopsis

Create a new page with a storage-format body. Pass --parent to create the
page below an existing page.

## Examples

confluence page create --space ENG --title "Runbook" --body-file body.html
  confluence page create --space ENG --title "Child" --parent 12345 --body-file body.html
  echo "<p>Hello</p>" | confluence page create --space ENG --title "Hello" --body-file -

## Usage

```text
confluence page create [flags]
```

## Options

```text
      --body string          Inline body string
      --body-file string     Path to body file, or '-' for stdin
      --body-format string   Body format: storage | wiki | view (default "storage")
      --parent string        Parent page id
      --space string         Space key (required)
      --title string         Page title (required)
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

- [confluence page](confluence_page.md)
