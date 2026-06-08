# confluence page screenshot

Capture a full-page browser screenshot via cdp

## Synopsis

Open a page in Chrome via cdp and capture a full-page screenshot. This is a
peer-tool wrapper; cdp must be installed and authenticated browser state must
already be available.

## Examples

confluence page screenshot 12345 --out verify.png
  confluence page screenshot 12345 --out verify.png --new-tab=false

## Usage

```text
confluence page screenshot <id> [flags]
```

## Options

```text
      --new-tab      Open a new tab instead of navigating an existing page (default true)
      --out string   Path to write screenshot image (required)
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
