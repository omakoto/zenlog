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

Zenlog-go is a complete write in Go, and no longer relies on script(1). 

## Quick start: install and setup (Bash)

To install, set up the Go SDK and run:

```
go install -u github.com/omakoto/zenlog-go/zenlog/cmd/zenlog 
```

Then add the following to your `.bashrc`.
```
. <(zenlog basic-bash-setup)
```

(Note this will overwrite `P0` and `PROMPT_COMMAND`, so if you don't like it, look at the script
and do whatever you want.)

### Using other shells

Any shell should work, as long as it supports some sort of "pre-exec" and "post-exec" hooks.

Look at the output of `zenlog basic-bash-setup` and figure it out.

## Using Zenlog

Once you setup your `.bashrc`, just run `zenlog` to start a new session.

As the [Caveats](#caveats) section describes, it's not recommended to use zenlog as a login shell.
Just change your terminal app's setting and use zenlog as a custom startup program.

Once in zenlog, output of all commands are stored under `$ZENLOG_DIR` (`$HOME/zenlog/` by default).

### Opening log files with commands

Try running `ls -l $HOME` and then `zenlog open-last-log`, which should invoke less
(`$ZENLOG_VIEWER`) with the output of the `ls` command.

### Opening log files with hotkeys

`basic-bash-setup` also sets up a hotkey `ALT+1` to open the last log.

Also if you have [A2H](https://github.com/omakoto/a2h-rs) installed, `ALT+2` will open the last log
in Google Chrome (change it with `$ZENLOG_RAW_VIEWER`) *with colors*.

## Below documentation is under construction...

[See also the readme of the old version.](https://github.com/omakoto/zenlog)
