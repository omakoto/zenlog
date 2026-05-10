[![Build Status](https://travis-ci.org/omakoto/zenlog.svg?branch=master)](https://travis-ci.org/omakoto/zenlog)

# Zenlog — automatic per-command shell logging

Zenlog wraps your login shell and silently saves every command's output to its
own log file, tagged with metadata (directory, git branch, timestamps, etc.).
No more `tee`, no more scrollback hunting.

## What you can do with it

| Want to… | How |
|---|---|
| Re-read the last command's output | Press `ALT+1` on the prompt (opens in `less`) |
| View it with full color in a browser | Press `ALT+2` (converts ANSI → HTML, requires `a2h`) |
| Pick any recent output interactively | Press `ALT+3` / `ALT+4` (uses `fzf`) |
| Insert a log filename into your command | Press `ALT+L` at any cursor position |
| Find URLs in recent output | `zenlog open-url` (picks with fzf, opens in browser) |
| Grep old output from `lsusb` | `grep … ~/zenlog/cmds/lsusb/SAN/` |
| Check what `make` printed last week | `less ~/zenlog/cmds/make/S` (symlink to latest) |
| Clean up old logs | `zenlog purge-log -p 30 -y` |

## Quick start

### 1. Install

```bash
mkdir -p "$HOME/src"
cd "$HOME/src"
git clone https://github.com/omakoto/zenlog.git
./zenlog/scripts/install.sh
export ZENLOG_SRC_DIR="$HOME/src/zenlog/"   # add this to your .*shrc too
```

### 2. Initialize

```bash
zenlog init
```

This creates `~/.zenlog.toml` and adds the setup hook to your `.bashrc` or
`.zshrc`.  The hook does **not** start zenlog automatically — it only wires up
the helper functions and hotkeys.

### 3. Start a session

```bash
zenlog
```

You're now in a logged shell.  Run any command, then press `ALT+1` to open its
output.  Type `exit` (or `exit 13` to restart zenlog) to leave the session.

> **Note (Bash):** Zenlog requires Bash 4.4+ for the `PS0` preexec hook.
> Run `bash --version` to check.

---

## Configuration: `~/.zenlog.toml`

Copy the sample from the source tree and edit to taste:

```bash
cp "$(zenlog zenlog-src-top)/dot_zenlog.toml" ~/.zenlog.toml
```

Key settings:

| Setting | Default | Description |
|---|---|---|
| `ZENLOG_DIR` | `$HOME/zenlog/` | Where log files are stored |
| `ZENLOG_START_COMMAND` | `exec $SHELL -l` | Shell started by `zenlog` |
| `ZENLOG_ALWAYS_NO_LOG_COMMANDS` | `vim`, `man`, `emacs`, `zenlog*`, … | Commands whose output is never saved |
| `ZENLOG_PREFIX_COMMANDS` | `sudo`, `time`, `command`, … | Prefixes stripped when determining the command name |
| `ZENLOG_AUTO_FLUSH` | `false` | Flush log files after every command |

Environment variables to set in `.bashrc` / `.zshrc`:

| Variable | Description |
|---|---|
| `ZENLOG_VIEWER` | Program to open SAN logs (default: `less`) |
| `ZENLOG_RAW_VIEWER` | Program to open RAW HTML logs (default: `google-chrome`) |
| `ZENLOG_BROWSER` | Browser for `zenlog open-url` (default: `google-chrome`) |
| `ZENLOG_NO_DEFAULT_PROMPT` | Set to `1` to disable the post-command status line |
| `ZENLOG_NO_DEFAULT_BINDING` | Set to `1` to disable the `ALT+*` hotkeys |

---

## Hotkeys (inside a zenlog session)

| Key | Action |
|---|---|
| `ALT+1` | Open last command output with `$ZENLOG_VIEWER` |
| `ALT+2` | Open last command output in browser with color (requires `a2h`) |
| `ALT+3` | Pick any recent log with fzf, open with `$ZENLOG_VIEWER` |
| `ALT+4` | Pick any recent log with fzf, open in browser with color |
| `ALT+L` | Insert last log filename at cursor (repeat to cycle older logs) |

---

## Log file structure

```
$ZENLOG_DIR/
├── SAN/YEAR/MM/DD/          # Sanitized (ANSI-stripped) — good for grep
├── RAW/YEAR/MM/DD/          # Raw output (with ANSI color codes)
├── ENV/YEAR/MM/DD/          # Metadata: pwd, git branch, exec time, …
│
├── cmds/
│   ├── cat/
│   │   ├── SAN/YEAR/…       # Sanitized logs for every "cat" invocation
│   │   ├── RAW/…
│   │   ├── ENV/…
│   │   ├── S                # Symlink → most recent sanitized log
│   │   ├── SS               # Symlink → second most recent
│   │   ├── R, RR, …         # Same for RAW
│   │   └── E, EE, …         # Same for ENV
│   ├── ls/
│   └── …
│
├── pids/PID/                # Per-session logs (same structure as cmds/)
├── tags/TAG/                # Per-tag logs (see "Log tagging" below)
│
├── S, SS, SSS, …            # Global most-recent sanitized logs (any shell)
├── R, RR, …                 # Global most-recent raw logs
└── E, EE, …                 # Global most-recent env logs
```

`S` always points to the most recent output, `SS` to the second most recent,
etc.  The `pids/` tree scopes these to a single shell session, so
`$ZENLOG_DIR/pids/$ZENLOG_PID/S` is your shell's last output even if another
shell ran a command in between.

### Log tagging

Add an inline comment to a command to create a named tag:

```bash
make -B    # full-build
```

Zenlog creates `$ZENLOG_DIR/tags/full-build/S` pointing to that run's output.

---

## Subcommands

### Session management

| Command | Description |
|---|---|
| `zenlog` | Start a new logged shell session |
| `zenlog init` | Interactive setup: create `~/.zenlog.toml` and patch shell RC |
| `zenlog pids` | List all active zenlog session PIDs |
| `zenlog flush` | Flush log buffers for the current session |
| `zenlog flush-all` | Flush log buffers for all sessions |
| `zenlog in-zenlog` | Exit 0 if running inside a zenlog session, 1 otherwise |

### Viewing logs

| Command | Description |
|---|---|
| `zenlog open-last-log [-r] [-e] [-p PID]` | Open previous command's log with `$ZENLOG_VIEWER` |
| `zenlog open-current-log [-r] [-e] [-p PID]` | Open current (most recent) command's log |
| `zenlog open-log` | Interactively pick a log with fzf and open it |
| `zenlog open-url` | Find URLs in recent logs, pick with fzf, open in browser |
| `zenlog cat-last-log` | Print previous command's log to stdout |
| `zenlog cat-last-log-content` | Same but strip the command-line header |
| `zenlog select-log [-r] [-e]` | Print a log filename selected via fzf |

### Finding log files

| Command | Description |
|---|---|
| `zenlog history [-r] [-e] [-n N] [-p PID]` | Print recent log filenames (up to 10) |
| `zenlog current-log [-r] [-e] [-p PID]` | Print filename of the most recent log |
| `zenlog last-log [-r] [-e] [-p PID]` | Print filename of the second most recent log |
| `zenlog all-commands [-r] [-e] [-n DAYS] [-c] [-l]` | List all logs with their command lines |

**`zenlog history` flags:**
- `-r` — show RAW log filename instead of SAN
- `-e` — show ENV log filename instead of SAN
- `-n N` — show only the Nth most recent log (1 = most recent)
- `-p PID` — use a specific zenlog session PID instead of `$ZENLOG_PID`

**`zenlog all-commands` flags:**
- `-n DAYS` — limit to logs within the last N days (default: 30)
- `-c` — limit to the current session only
- `-l` — print filenames only, no command lines

### Maintenance

| Command | Description |
|---|---|
| `zenlog purge-log -p DAYS [-y] [-P] [-b]` | Delete logs older than DAYS days |
| `zenlog du [du-options]` | Show disk usage of the log directory |
| `zenlog free-space` | Print free bytes on the log filesystem |
| `zenlog ensure-log-dir` | Create the log directory if it doesn't exist |
| `zenlog update` | Download and compile the latest zenlog from GitHub |

**`zenlog purge-log` flags:**
- `-p DAYS` — (required) delete logs older than this many days
- `-y` — skip the confirmation prompt
- `-P` — dry run: print files that would be deleted instead of deleting them
- `-b` — run in the background

```bash
# Examples
zenlog purge-log -p 30          # Remove logs older than 30 days (asks for confirmation)
zenlog purge-log -p 90 -y       # Force-remove logs older than 90 days
zenlog purge-log -p 7 -P        # Dry run: show what would be deleted
```

### Shell setup (advanced / manual)

| Command | Description |
|---|---|
| `zenlog basic-bash-setup` | Output the Bash setup snippet (source into `.bashrc`) |
| `zenlog basic-zsh-setup` | Output the Zsh setup snippet (source into `.zshrc`) |
| `zenlog sh-helper` | Output common helper shell functions |
| `zenlog zenlog-src-top` | Print the zenlog source directory |
| `zenlog zenlog-bin` | Print the path to the zenlog binary |

### Scripting / internals

These are used by the shell hooks and scripts; you rarely need them directly.

| Command | Description |
|---|---|
| `zenlog start-command [-e ENV] CMDLINE…` | Signal the logger that a command is starting |
| `zenlog end-command [-n]` | Signal the logger that the command finished |
| `zenlog write-to-logger` | Pipe stdin into the logger (counted as output) |
| `zenlog write-to-outer` | Pipe stdin to the terminal, bypassing the logger |
| `zenlog outer-tty` | Print the outer TTY device filename |
| `zenlog logger-pipe` | Print the named pipe path to the logger |
| `zenlog check-bin-update` | Print a message if the binary was updated since session start |
| `zenlog list-logs [DIR]` | List the log directory tree (F=file, L=symlink) |

---

## Manual shell setup

If `zenlog init` doesn't cover your shell, add the setup line by hand.

**Bash** (add to `~/.bashrc`):
```bash
. <(zenlog basic-bash-setup)
```

**Zsh** (add to `~/.zshrc`):
```zsh
. <(zenlog basic-zsh-setup)
```

For any POSIX-compatible shell with preexec/postexec hooks, look at the output
of `zenlog basic-bash-setup` and adapt the hook wiring to your shell.

---

## How it works

Zenlog creates a PTY, starts your shell inside it, and tees all I/O through its
own logging goroutines.  Each command's output is separated by the `PS0` (Bash)
or `preexec` (Zsh) and `PROMPT_COMMAND` / `precmd` hooks.

Output is written to three parallel trees:

- **SAN** — ANSI escape sequences stripped (grep-friendly)
- **RAW** — byte-for-byte original output (color-accurate)
- **ENV** — metadata: working directory, git branch, exit status, timestamps

Symbolic links in `cmds/`, `pids/`, and `tags/` provide O(1) access to the
most recent N outputs per command, session, and tag.

[Legacy version (Perl/Ruby)](https://github.com/omakoto/zenlog-legacy)
