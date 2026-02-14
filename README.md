# Storage

A CLI and terminal-based Google Drive file browser built with [Cobra](https://github.com/spf13/cobra) and [Bubble Tea](https://github.com/charmbracelet/bubbletea).

## Prerequisites

1. **Go** (1.25.6+) must be installed
2. **Google Cloud OAuth2 credentials** are required

## Setup Google Credentials

1. Go to the [Google Cloud Console Credentials page](https://console.cloud.google.com/apis/credentials)
2. Create an **OAuth 2.0 Client ID** (type: "Desktop app")
3. Enable the **Google Drive API** in your Google Cloud project
4. Download the credentials JSON file
5. Save it to `~/.config/storage/credentials.json`

## Usage

```bash
# Build the binary
go build -o storage ./cmd/storage/

# Or run directly with go run
go run ./cmd/storage/ <command>
```

### Commands

```bash
# Launch the interactive TUI file browser
storage tui

# List files in your Drive root
storage list

# List files in a specific folder
storage list <folder-id>

# Show help
storage --help
```

## First-Time Authentication

On the first run, the app will:

1. Print an authorization URL to the terminal
2. You open that URL in your browser and grant Google Drive read-only access
3. The browser redirects to `http://127.0.0.1:9874/callback`
4. A token is cached at `~/.config/storage/token.json` for future runs

## TUI Keyboard Controls

| Key                | Action           |
| ------------------ | ---------------- |
| `j` / `Down`       | Move cursor down |
| `k` / `Up`         | Move cursor up   |
| `Enter` / `l`      | Enter a folder   |
| `Backspace` / `h`  | Go back to parent|
| `q` / `Ctrl+C`     | Quit             |

## Project Structure

```
storage/
├── cmd/storage/            # Entry point
├── pkg/
│   ├── core/               # Domain types, Backend interface, FileService
│   ├── storage/gdrive/     # Google Drive backend implementation
│   └── auth/               # Google OAuth2 authentication
├── internal/
│   ├── cli/                # Cobra commands (root, list, tui)
│   └── tui/                # Bubble Tea interactive file browser
```

## License

MIT
