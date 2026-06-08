# confluence page

Pages (view, search, children, direct-children, descendants, ancestors, history, versions, create, update, publish, delete, purge, url, screenshot)

## Synopsis

Page operations.

## Examples

confluence page view 12345 --markdown
  confluence page search --cql "type=page AND space=ENG"
  confluence page children 12345 --recursive
  confluence page direct-children 12345 --json
  confluence page descendants 12345 --depth 2
  confluence page create --space ENG --title "Runbook" --body-file body.html
  confluence page delete 12345 --force
  confluence page url 12345
  confluence page screenshot 12345 --out verify.png

## Usage

```text
confluence page
```

## Commands

- `ancestors` - List ancestor pages of a content id
- `children` - List children of a page
- `create` - Create a new page
- `delete` - Move a page to trash
- `descendants` - List descendant pages/content under a page
- `direct-children` - List direct mixed children of a page
- `history` - Fetch page history
- `publish` - Upload attachments, then update page body
- `purge` - Permanently delete a trashed page
- `screenshot` - Capture a full-page browser screenshot via cdp
- `search` - CQL-powered page search
- `update` - Update an existing page (title and/or body)
- `url` - Print the browser URL for a page
- `versions` - List version records for a page
- `view` - Fetch a page (or other content) by id

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
