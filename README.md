# zsh-dcmd

Directory-Specific Command Manager

A simple plugin to create, list, and run project-specific shell commands with variable support. Designed for Zsh, for personal use. Rewritten in Go for speed and maintainability.

---

## Features

#### Context-Aware Commands
Define commands that are only available when you `cd` into a specific project directory. Stop remembering complex build, test, or deploy commands for each project.

#### Global Commands
A separate `dcmd global` list for commands you need everywhere, regardless of your current directory.

#### Variables
Define commands with `{{variable}}` placeholders. When you run the command, the plugin will automatically prompt you to fill in the values, making your scripts reusable.

#### Non-Intrusive by Design

- All command configurations are stored centrally in your home directory (~/.config/dcmd).
- Your project directories are never modified. This means you don't have to add config files to your `.gitignore` or pollute your project with personal-use files. This keeps your shortcuts independent of the codebase.

#### Maintenance

If you move or rename a project folder, your commands aren't lost. `dcmd link` helps you find these "orphaned" configurations and re-link them to the project's new location.

Easily `delete` commands or `renumber` your entire list to keep it clean and organized.

#### Safe Execution

Commands are not run immediately. Instead, they are placed onto your command line for review or edit before you press Enter.

---

## Platform Availability

| Platform | Status | Notes |
|----------|--------|-------|
| macOS    | Supported | Primary development platform. `dcmd config` opens files with `open`. |
| Linux    | Supported | `dcmd config` opens files with `xdg-open`. |
| Windows  | Partial | The Go binary works, but the Zsh wrapper requires a Zsh-compatible shell (e.g. WSL, Git Bash). |

The Go binary itself is cross-platform. The Zsh wrapper (`dcmd.plugin.zsh`) is what integrates it into your shell — it requires Zsh and the `print -z` builtin for safe command buffering.

---

## Installation

No Go required — the plugin automatically downloads a pre-built binary on first run.

### Manual (Generic Zsh)

Clone this repository somewhere on your machine, for example in `~/.zsh/zsh-dcmd`:

```zsh
git clone https://github.com/karpikorop/zsh-dcmd-go.git ~/.zsh/zsh-dcmd
```

Add the following to your `.zshrc`:

```zsh
source ~/.zsh/zsh-dcmd/dcmd.plugin.zsh
```

### Oh My Zsh

Clone this repository into your Oh My Zsh custom plugins directory:

```zsh
git clone https://github.com/karpikorop/zsh-dcmd-go.git ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/dcmd
```

Add dcmd to the list of plugins in your `~/.zshrc` file:

```zsh
plugins=(
# ... other plugins
dcmd
)
```

Restart your shell (or run `source ~/.zshrc`) to apply the changes.

The first time you run `dcmd`, the wrapper will download the correct binary for your platform from [GitHub Releases](https://github.com/karpikorop/zsh-dcmd-go/releases). If the download fails (unsupported platform, no internet), it will fall back to building from source if [Go](https://go.dev/dl/) is available.

### Updating

To update to the latest version, pull the latest changes and remove the old binary:

```zsh
cd ~/.zsh/zsh-dcmd && git pull && rm -f dcmd
```

The next time you run `dcmd`, the wrapper will download the latest release.

---

## Configuration

By default, `dcmd` stores all its configuration files in `~/.config/dcmd`.

If you want to change this (e.g., to sync with your dotfiles), you can set the `DCMD_CONFIG_DIR` variable in your `.zshrc` before the plugin is sourced.

```zsh
export DCMD_CONFIG_DIR="${HOME}/.my-dotfiles/dcmd"
```

---

## Usage

dcmd is designed to be run from within your project directories. It also provides a global sub-command for commands you need everywhere.

### Project Commands (Default)

These commands act on the project associated with your current directory.

- `dcmd` Lists all available commands for the current project.

- `dcmd <number>` Selects a command from the project list to run. The command is placed on your command line for review before you press Enter. If the command has {{variables}}, you will be prompted to fill them in first.

- `dcmd add` Starts an interactive prompt to add a new command to the current project.

- `dcmd add --all` Opens an interactive menu to add a new command to multiple selected projects (or all of them) at once.

- `dcmd delete` Shows the current project's commands and prompts you to select one to delete. Automatically renumbers the remaining commands.

- `dcmd config [project]` Opens the config file for the current project (or a specified project) in your default GUI editor.

- `dcmd renumber` Re-numbers all commands in the current project's list sequentially from 1. This is useful for cleaning up the list.

- `dcmd link` Helps find and re-link an "orphaned" config file to the current directory. Use this if you have moved or renamed a project folder and want to recover its old dcmd list.

### Global Commands

These commands are available everywhere, regardless of your current directory.

- `dcmd global` Lists all available global commands.

- `dcmd global <number>` Selects a command from the global list to run, placing it on your command line for review.

- `dcmd global add` Starts an interactive prompt to add a new global command.

### Help

- `dcmd help` Displays the help message with a list of all available commands and subcommands.

- `dcmd help <command>` Shows detailed help for a specific command (e.g. `dcmd help global`).
