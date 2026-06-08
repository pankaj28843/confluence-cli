# confluence attachment replace

Create or replace an attachment on a page

## Synopsis

Upload an attachment. If the filename already exists on the page, the
existing attachment data is updated as a new version.

## Examples

confluence attachment replace --page 12345 --file ./report.pdf
  confluence attachment replace --page 12345 --file ./logo.png --comment "v2 logo"
  cat report.pdf | confluence attachment replace --page 12345 --file - --file-name report.pdf

## Usage

```text
confluence attachment replace [flags]
```

## Options

```text
      --comment string     Optional attachment comment
      --file string        Local file path, or '-' for stdin (required)
      --file-name string   Filename to use when --file is '-'
      --page string        Content id to attach to (required)
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

- [confluence attachment](confluence_attachment.md)
