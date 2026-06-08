# confluence docs markdown

Generate Markdown CLI reference files

## Synopsis

Generate one timestamp-free Markdown file per command. The output is stable
and suitable for docs sites, search indexes, and LLM context.

## Examples

confluence docs markdown --out docs/cli
  confluence docs markdown --out /tmp/confluence-cli-docs

## Usage

```text
confluence docs markdown [flags]
```

## Options

```text
      --out string   Output directory (required)
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

- [confluence docs](confluence_docs.md)
