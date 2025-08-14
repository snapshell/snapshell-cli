# SnapShell CLI

Convert CLI output to clean, styled, and shareable web snapshots.

## Features
- Converts raw CLI output (e.g., `terraform plan`, `npm audit`) into web snapshots
- Auto-detects output type for labeling
- Optional authentication for private snapshots
- Supports reading from files or stdin
- Expiration and privacy controls

## Installation

### From GitHub Releases

1. Go to the [Releases page](https://github.com/snapshell/snapshell-cli/releases).
2. Download the appropriate archive for your OS and architecture:
   - Linux/macOS: `snapshell_<version>_linux_amd64.tar.gz` or `snapshell_<version>_darwin_arm64.tar.gz`
   - Windows: `snapshell_<version>_windows_amd64.zip`
3. Extract the archive:
   - Linux/macOS:
     ```sh
     tar -xzf snapshell_<version>_<os>_<arch>.tar.gz
     sudo install -m 755 snapshell /usr/local/bin/
     ```
   - Windows:
     Unzip the file and place `snapshell.exe` somewhere in your PATH (e.g., `C:\Windows\System32` or add a custom directory to your PATH).

### Build from Source

```bash
git clone https://github.com/snapshell/snapshell-cli.git
cd snapshell-cli
go build -o snapshell .
```

Or use Go:

```bash
go install github.com/snapshell/snapshell-cli@latest
```

## Usage

### Pipe CLI Output

```bash
terraform plan | snapshell --label="My Plan"
npm audit | snapshell --label="Security Audit"
```

### Read from File

```sh
snapshell --file=plan.txt --label="My Plan"
```

### Auto-labeling
- If `--label` is omitted and using `--file`, the file name is used as the label.
- If `--label` is omitted and using stdin, the label is set to the detected type plus a timestamp (e.g., `terraform-2025-08-14_13-45-00`).

### Authentication

Login via browser to enable private snapshots:

```sh
snapshell login
```

This will open a browser for you to authenticate with https://app.snapshell.dev

Logout:

```sh
snapshell logout
```

## Flags
- `--label`        Snapshot label (optional)
- `--type`         Snapshot type (auto-detected if not specified)
- `--private`      Make snapshot private (default: true)
- `--expires`      Snapshot expiration in days (default: 30)
- `--file`         Read from file instead of stdin
- `--api`          API base URL (default: https://app.snapshell.dev)

## Example Output

After running, SnapShell prints a shareable URL to your snapshot:

```
https://app.snapshell.dev/snapshots/abc123
```

## Development

- Main logic in `main.go`, commands in `pkg/commands`, auth in `pkg/auth`, snapshot logic in `pkg/snapshot`
- Uses Cobra for CLI

## License

MIT
