[![Build Status](https://travis-ci.org/omakoto/zenlog-go.svg?branch=master)](https://travis-ci.org/omakoto/zenlog-go)

# Zenlog -- no more tee(1)-ing.

Zenlog wraps a login shell and automatically saves all the output of each command to a separate log
file, along with metadata such as each command start/finish time, the current directly, etc.

You can now keep all the output of every single command you'll execute as long as you have extra
disk space.

If you use a command line editor such as Emacs or Vi, you can blacklist them to prevent saving
output from such commands.

The primary target shell is Bash 4.4 or later (4.4 is required because you need to tell when each
command starts with the `P0` prompt), but any shell with similar syntax an a pre-exec hook should
work with zenlog.

# How it works

[See the readme of the old version.](https://github.com/omakoto/zenlog)

Zenlog-go uses the same idea, but is a complete write in Go, and no longer relies on script(1).
Instead it'll create a PTY by itself. 

## Quick start: Install and set up

To install, set up the Go SDK and run:

```
go get -v -u github.com/omakoto/zenlog-go/zenlog/cmd/zenlog 
```

Then, create `~/.zenlog.toml` and update `.bashrc` (or `.zshrc`) by running:

```
zenlog init
```

Then, start a new zenlog session by running:

```
zenlog
```

Then, try running `ls -l` and then press `ALT+1`. The output of the `ls` command should open
in `less`. (If the hotkey doesn't work, then run `zenlog open-last-log` instead.)


## Manual bash Setup

`zenlog init` will add `. <(zenlog basic-bash-setup)`, which overwrites `P0` and `PROMPT_COMMAND`.

```
. <(zenlog sh-setup)
```
 
TO BE WRITTEN

## Manual zsh setup (experimentawl)

TO BE WRITTEN

### Using other shells

Any shell should work, as long as it supports some sort of "pre-exec" and "post-exec" hooks.

Look at the output of `zenlog basic-bash-setup` and figure it out.


[See also the readme of the old version.](https://github.com/omakoto/zenlog)
