# Contributing to confluence-cli

Thanks for taking the time to contribute! A few ground rules keep
reviews fast.

## Dev setup

```bash
git clone https://github.com/pankaj28843/confluence-cli
cd confluence-cli
go mod download
make setup     # installs the pre-commit hook
```

Requirements: Go 1.25+, `jq` on `$PATH` for `--jq` tests.

## Running the test suite

```bash
make test      # fast — go test ./... -count=1 -timeout 60s
make race      # go test ./... -race
make coverage  # internal/ coverage to COVERAGE.txt
```

## House style

See `AGENTS.md`. Short version:

- `--json` / `--jq` / `--template` / `--timing` / `--debug` on every
  command (they are persistent root flags).
- Every verb's `Long:` ends with an `Examples:` block.
- Errors that the user can fix by changing env vars / tokens return via
  `newConfigError(err)` → exit 2.
- One runtime dep (cobra) + one converter dep (goquery). If you need
  another dep, open an issue first.
- All example URLs in docs use `https://wiki.example.com` (Server/DC) or
  `https://example.atlassian.net/wiki` (Cloud). No real hostnames.

## Pull requests

1. Work on a feature branch.
2. Run `make race` and `make coverage` before pushing.
3. Explain the *why* in the PR description, not just the *what*.
4. If you change the REST surface, add an httptest unit test.

## Reporting bugs

Please include:

- Your Confluence flavor (Server/DC vs Cloud) + version.
- The exact command invocation (redact tokens).
- `confluence --debug <...> 2>&1` output (Authorization is auto-redacted).
