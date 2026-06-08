# confluence attachment download

Download an attachment by name

## Synopsis

Download an attachment. --output file writes to disk; omit for stdout.

## Examples

confluence attachment download --page 12345 --name logo.png --output ./logo.png
  confluence attachment download --page 12345 --name report.pdf > report.pdf

## Usage

```text
confluence attachment download [flags]
```

## Options

```text
      --name string     Attachment file name (required)
      --output string   Write to file (default: stdout)
      --page string     Content id (required)
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
