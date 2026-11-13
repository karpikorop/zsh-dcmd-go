# zsh-dcmd

Directory-Specific Command Manager

A simple tool to create, list, and run project-specific shell commands with variable support. Designed for Zsh, for personal use. Written in Go for speed and maintainability.

---

## Features

#### Context-Aware Commands
Define commands that are only available when you `cd` into a specific project directory or any of its subfolders. Stop remembering complex build, test, or deploy commands for each project.

One project can have several folders (for example clones or worktrees), so they share the same commands.

#### Global Commands
A separate `dcmd global` list for commands you need everywhere, regardless of your current directory.

#### Variables
Define commands with `{{variable}}` placeholders. When you run the command, dcmd prompts you to fill in the values, making your scripts reusable.

- `{{name:default}}` gives a placeholder a default value. The prompt starts with it, so you only press Enter.
- Project variables fill a placeholder automatically. Set one with `dcmd var set project my-app-dev`, and every `{{project}}` in that project uses it without a prompt.

#### Non-Intrusive by Design

- All commands are stored centrally in your home directory (`~/.config/dcmd`), one plain TOML file per project.
- Your project directories are never modified. This means you don't have to add config files to your `.gitignore` or pollute your project with personal-use files. This keeps your shortcuts independent of the codebase.

#### Maintenance

If you move or rename a project folder, your commands aren't lost. `dcmd link` helps you find these "orphaned" projects and re-link them to the project's new location.

Command numbers stay the same when you delete a command. Run `renumber` when you want to close the gaps. Use `copy` to reuse commands from another project.

#### Safe Execution

Commands are not run immediately. Instead, they are placed onto your command line for review or edit before you press Enter.

---

## Platform Availability

| Platform | Status | Notes |
|----------|--------|-------|
| macOS    | Supported | Primary development platform. `dcmd config` opens files with `open`. |
| Linux    | Supported | `dcmd config` opens files with `xdg-open`. |
| Windows  | Partial | The Go binary works, but the Zsh wrapper requires a Zsh-compatible shell (e.g. WSL, Git Bash). |

The Go binary itself is cross-platform. The Zsh wrapper (printed by `dcmd init zsh`) is what integrates it into your shell — it requires Zsh and the `print -z` builtin for safe command buffering.

---

## Installation

Install with [Homebrew](https://brew.sh) (macOS and Linux):

```zsh
brew install karpikorop/tap/dcmd
```

Add the following to your `~/.zshrc`:

```zsh
eval "$(dcmd init zsh)"
```

Restart your shell (or run `source ~/.zshrc`) to apply the changes.

### Updating

```zsh
brew upgrade dcmd
```

---

## Configuration

By default, `dcmd` stores its files in `~/.config/dcmd`:

- `projects/<name>.toml` — the folders, variables and commands of one project
- `global.toml` — the global commands

If you want to change this (e.g., to sync with your dotfiles), set the `DCMD_CONFIG_DIR` variable in your `.zshrc`:

```zsh
export DCMD_CONFIG_DIR="${HOME}/.my-dotfiles/dcmd"
```

---

## Usage

dcmd is designed to be run from within your project directories. It also provides a global sub-command for commands you need everywhere.

### Project Commands (Default)

These commands act on the project associated with your current directory.

- `dcmd` Lists all available commands for the current project.

- `dcmd <number>` Selects a command from the project list to run. The command is placed on your command line for review before you press Enter. If the command has {{variables}}, you will be prompted to fill them in first. Pass `--var NAME=VALUE` to skip a prompt.

- `dcmd show <number>` Shows one command with its template and variables.

- `dcmd add` Starts an interactive prompt to add a new command to the current project.

- `dcmd add --all` Opens an interactive menu to add a new command to multiple selected projects (or all of them) at once.

- `dcmd edit [number]` Edits the command and the description of a saved command.

- `dcmd delete [number]` Deletes a command. Without a number, it shows the list and asks you to select one. The other commands keep their numbers.

- `dcmd renumber` Re-numbers all commands in the current project's list sequentially from 1. This is useful for cleaning up the list.

- `dcmd copy <project>` Copies commands from another project. Use `--ids 1,4,7` to copy only some of them. Commands that already exist are skipped.

- `dcmd var` Lists the project variables. `dcmd var set <name> <value>` and `dcmd var unset <name>` change them.

- `dcmd projects` Lists all projects with their folders.

- `dcmd path add <project>` Adds the current folder to an existing project. Use this for clones and worktrees.

- `dcmd link` Helps find and re-link an "orphaned" project to the current directory. Use this if you have moved or renamed a project folder and want to recover its old dcmd list.

- `dcmd config [project]` Opens the file of the current project (or a specified project) in your default GUI editor.

### Global Commands

These commands are available everywhere, regardless of your current directory.

- `dcmd global` Lists all available global commands.

- `dcmd global <number>` Selects a command from the global list to run, placing it on your command line for review.

- `dcmd global add` Starts an interactive prompt to add a new global command.

- `dcmd global edit` Edits a global command.

Most project commands also accept `--global` to work on the global list instead, for example `dcmd delete 3 --global`.

### Help

- `dcmd help` Displays the help message with a list of all available commands and subcommands.

- `dcmd help <command>` Shows detailed help for a specific command (e.g. `dcmd help global`).

---

## Use with AI Agents

Scripts and AI agents can use dcmd too. These flags skip all prompts:

| Flag | Works with | What it does |
|------|------------|--------------|
| `--json` | `dcmd`, `dcmd global`, `show`, `var`, `projects` | Prints the result as JSON, with the variables and defaults of each command |
| `--var NAME=VALUE` | `dcmd <number>`, `dcmd global <number>` | Gives a value for a variable, so dcmd does not ask for it |
| `--print` | `dcmd <number>`, `dcmd global <number>` | Prints the final command to stdout instead of your prompt |
| `--command`, `--description` | `add`, `edit` | Sets the command and the description |
| `--confirm` | `add`, `edit` | Marks a command (for example a deploy) that agents must ask you about before they run it |
| `--yes` | `delete` | Deletes without asking |

```zsh
dcmd --json
dcmd 4 --var FunctionName=invoca --print
dcmd add --command "yarn build" --description "Build"
```
