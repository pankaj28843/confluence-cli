# confluence

Atlassian Confluence CLI (Server/DC + Cloud)

## Synopsis

confluence - Atlassian Confluence CLI
Version: dev (built unknown, commit unknown)

Query and act on Atlassian Confluence from the terminal. Supports both
Confluence Server / Data Center (Bearer PAT, /rest/api/...) and Confluence
Cloud (Basic email:token, /rest/api/ + /api/v2/ relative to a
CONFLUENCE_URL that already ends in /wiki).

## Requires environment variables

CONFLUENCE_URL                Base URL, e.g. https://wiki.example.com
                                (Server/DC) or https://example.atlassian.net/wiki (Cloud)
  Server/DC:
    CONFLUENCE_PAT              Personal access token (alias:
                                CONFLUENCE_PERSONAL_ACCESS_TOKEN)
  Cloud:
    CONFLUENCE_EMAIL            Atlassian account email
    CONFLUENCE_API_TOKEN        API token from id.atlassian.com
  CONFLUENCE_FLAVOR             Optional override: server|cloud
  CONFLUENCE_DEFAULT_SPACE      Optional default for --space

## Workflow

confluence doctor                                         Health + auth + flavor probe
  confluence space list                                     Spaces in the site
  confluence page search --cql "type=page AND space=ENG"    CQL-powered page search
  confluence blogpost list --space ENG --limit 10           Blog posts in a space
  confluence page view 12345 --markdown                     Fetch + render a page
  confluence page children 12345 --recursive --json         Page-tree traversal
  confluence comment list --page 12345                      Inline + footer + resolved
  confluence label add --page 12345 --label needs-review    Low-risk mutation
  confluence property content get --page 12345 --key foo    Content property lookup
  confluence task list --status incomplete --limit 10       Cloud content tasks
  confluence operation list --page 12345                    Permitted operations
  confluence like count --page 12345                         Cloud like count
  confluence body convert --to view --value @body.xml       Convert storage/editor bodies
  confluence database view 12345 --json                     Cloud database details
  confluence folder children 12345 --json                   Cloud content-tree children
  confluence custom-content list --type ac:example --json   Cloud custom content
  confluence macro body --page 12345 --version 2 --macro-id id  Fetch macro body
  confluence template list --limit 10                       Cloud content templates
  confluence permission space list --space ENG --json       Space permission assignments
  confluence search all "release process" --json            Unified fan-out (code/space/user/attachment)
  confluence api /rest/api/user/current                     Raw REST passthrough

## Usage

```text
confluence
```

## Commands

- `api` - Call any Confluence REST endpoint (escape hatch)
- `attachment` - Attachments (list, versions, version, download, upload, replace, delete)
- `blogpost` - Blog posts (list, view, versions, version, create, update, delete, purge)
- `body` - Content body utilities
- `comment` - Comments (list/versions/version/add/update/delete)
- `custom-content` - Cloud custom content reads
- `database` - Database reads
- `docs` - Generate CLI reference documentation
- `doctor` - Verify environment, auth, and flavor detection
- `folder` - Folder reads
- `group` - Groups (list, view, picker, members, hierarchy)
- `label` - Content and space labels
- `like` - Cloud likes (count, users)
- `macro` - Macro body utilities
- `operation` - Permitted operations (list)
- `page` - Pages (view, search, children, direct-children, descendants, ancestors, history, versions, version, create, update, publish, delete, purge, url, screenshot)
- `permission` - Space permission reads
- `property` - Content and space properties
- `restriction` - Content restrictions (list)
- `search` - Search content, spaces, users, attachments, or all
- `smart-link` - Smart Link reads
- `space` - Spaces (list, view)
- `task` - Tasks (list, view, complete, reopen, long)
- `template` - Cloud content and blueprint templates
- `user` - Users (current, view, search, bulk)
- `version` - Print version, build time, and commit
- `watcher` - Watchers (content, space, status)
- `whiteboard` - Whiteboard reads

## Options

```text
      --debug             Log HTTP requests to stderr (Authorization header redacted)
      --jq string         Filter JSON output through a jq expression (requires jq on PATH)
      --json              Output as JSON (machine-readable)
      --template string   Render JSON output through a Go text/template
      --timing            Show execution time on stderr
```
