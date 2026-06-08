# confluence page children

List children of a page

## Synopsis

List child content under a parent id. Default childType is 'page'.
--recursive traverses breadth-first.

## Examples

confluence page children 12345
  confluence page children 12345 --recursive --json
  confluence page children 12345 --type comment

## Usage

```text
confluence page children <id> [flags]
```

## Options

```text
      --limit int     Max children per parent (hard cap 200) (default 50)
      --recursive     Walk the full page tree
      --type string   Child type: page | comment | attachment (default "page")
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
