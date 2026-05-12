# Windows developer setup

This guide is for developers who already have a `confluence-windows-amd64.exe` binary available. It does not describe a distribution channel, installer, signing flow, or public release process.

## Recommended layout

Use a per-user tools directory and rename the platform binary to `confluence.exe` so shells and coding agents can call `confluence` without the OS/architecture suffix:

```powershell
mkdir "$env:USERPROFILE\tools" -Force
copy .\confluence-windows-amd64.exe "$env:USERPROFILE\tools\confluence.exe"
```

Add the tools folder to your user-level PATH:

```powershell
$tools = "$env:USERPROFILE\tools"
$p = [Environment]::GetEnvironmentVariable("PATH", "User")
if (($p -split ';') -notcontains $tools) {
    [Environment]::SetEnvironmentVariable("PATH", "$tools;$p", "User")
}
```

Open a new PowerShell, terminal, Claude Code session, or Copilot CLI session after changing PATH. Already-open processes do not inherit the updated user environment.

## Build the developer binary

From this repository, a developer can cross-compile the Windows binary with:

```bash
make windows
```

The output is `build/confluence-windows-amd64.exe`. Copy that file to the Windows machine, then follow the layout above.

## Sanity check

In a new PowerShell session:

```powershell
Get-Command confluence
confluence version
```

If `Get-Command` cannot find `confluence`, confirm the file exists and that the new terminal picked up the user PATH change:

```powershell
Test-Path "$env:USERPROFILE\tools\confluence.exe"
[Environment]::GetEnvironmentVariable("PATH", "User")
```

## Configure Confluence

Set user-level environment variables so future terminals and coding-agent sessions inherit the same configuration.

### Server / Data Center (Bearer PAT)

```powershell
[Environment]::SetEnvironmentVariable("CONFLUENCE_URL", "https://wiki.example.com", "User")
[Environment]::SetEnvironmentVariable("CONFLUENCE_PAT", "<your-pat>", "User")
```

### Cloud (Basic email + API token)

```powershell
[Environment]::SetEnvironmentVariable("CONFLUENCE_URL", "https://example.atlassian.net/wiki", "User")
[Environment]::SetEnvironmentVariable("CONFLUENCE_EMAIL", "you@example.com", "User")
[Environment]::SetEnvironmentVariable("CONFLUENCE_API_TOKEN", "<your-api-token>", "User")
```

### Optional

```powershell
[Environment]::SetEnvironmentVariable("CONFLUENCE_FLAVOR", "server", "User")
[Environment]::SetEnvironmentVariable("CONFLUENCE_DEFAULT_SPACE", "ENG", "User")
[Environment]::SetEnvironmentVariable("SSL_CERT_FILE", "C:\\path\\to\\ca.pem", "User")
```

Open a new terminal after changing user-level environment variables.

## Validate configuration

```powershell
confluence doctor
```

Expected healthy output includes the detected flavor, base URL, authenticated user, and `OK -- confluence is ready to use.` Exit code `2` means user-fixable configuration such as a missing URL or token; exit code `1` means an unexpected error.

## Troubleshooting

```powershell
Get-Command confluence
confluence version
confluence doctor
Get-ChildItem Env:CONFLUENCE*
Get-Item Env:SSL_CERT_FILE -ErrorAction SilentlyContinue
```

Common fixes:

- Reopen the terminal after PATH or `CONFLUENCE_*` changes.
- Check that `%USERPROFILE%\tools\confluence.exe` exists.
- Use `CONFLUENCE_URL=https://wiki.example.com` for Server/Data Center style URLs and `CONFLUENCE_URL=https://example.atlassian.net/wiki` for Cloud style URLs.
- Set `SSL_CERT_FILE` only when the Confluence host uses a private CA bundle.
- Keep tokens out of scripts, screenshots, shell history snippets, and issue reports.

## Notes for maintainers

- This folder intentionally contains documentation only. Do not check in generated `.exe` files.
- Do not add installers, MSI packages, signing steps, publishing commands, or scheduled tasks here.
- Do not add VBScript wrappers unless the CLI grows a real background workflow that requires hidden Windows Script Host execution.
