[![Build Status](https://travis-ci.org/omakoto/zenlog.svg?branch=master)](https://travis-ci.org/omakoto/zenlog)

# Zenlog -- no more tee-ing

~~2020-02-22 NOTE: Looks like Go stopped putting the source file under `$HOME/go/src`.~~
~~Without it, the external subcommands like `zenlog init`, etc will not work~~
~~To make it work again, follow the new install instructions below...~~

**2025-01-05 NOTE: I tried fixing the above problem. Now you should be able to install zenlog with `go install`.**
**(`git clone` shouldn't be needed.)**

Zenlog wraps a login shell and automatically saves all the output of each command to a separate log
file, along with metadata such as each command start/finish time, the current directory, etc.

It has various applications:
 - Want to open the output of the last command in my faviroite editor!
    - With the default installation you can do it with simply pressing `ALT+1` on the command line.
     
 - "What was the output of `lsusb` command I ran a month ago?"
    - Zenlog keeps log files in such a way that it's easy to find out a specific output, with
        meta-information such as the directory it was executed in, the git branch, the execution
        time, etc.
        
 - The previous command output had an HTTP link. I want to open it in the browser.
   - Zenlog provides various commands to access the output of previous commands, to help write a script like this. 
    

The primary target shell is Bash 4.4 or later (Zenlog requires the `PS0` (aka preexec) hook added
in Bash 4.4), but any shell with similar syntax with a pre-exec hook should work with zenlog. It
comes with an installation script that supports both Bash and Zsh.

# How it works

[See the readme of the old version.](https://github.com/omakoto/zenlog-legacy)

Zenlog uses the same idea as the previous perl/ruby versions,
but is a complete write in Go, and no longer relies on script(1), and 
instead it'll create a PTY by itself. 

## Quick start: Install and set up

To install, set up the Go SDK and follow the below instructions:

The following command __should__ be enough to install zenlog. (`git clone` shouldn't be needed.)

**NOTE: `-a` is needed to ensure all source files are rebuilt, which is needed to detect source directory**
```
go install -a github.com/omakoto/zenlog/zenlog/cmd/zenlog@latest
```


Then, run the following command to create `~/.zenlog.toml` and update `.bashrc` (or `.zshrc`):

```
zenlog init
```

Then, start a new zenlog session by running:

```
zenlog
```

Then, try running `ls -l` and then press `ALT+1`. The output of the `ls` command should open
in `less`. (If the hotkey doesn't work, then run `zenlog open-last-log` instead.)

## Customization

 - [`~/.zenlog.toml`](dot_zenlog.toml) contains various configuration such as the log directory. 

 - Set `$ZENLOG_VIEWER` and `$ZENLOG_RAW_VIEWER` to change what command to use to open log files.
Set it in `.bashrc` / `.zshrc`.

## Manual Bash / Zsh Setup

 - Create `.zenlog.toml` in your home directory:
 
```bash
cp "$(zenlog zenlog-src-top)/dot_zenlog.toml" "$HOME/.zenlog.toml"
``` 

 - Then, if you're using Bash, add the following line to your `~/.bashrc`.
```bash
. <(zenlog basic-bash-setup)
```

 - Then, if you're using Zsh, add it to your `~/.zrc`.
```zsh
. <(zenlog basic-zsh-setup)
```

### Using other shells

Any shell should work, as long as it supports some sort of "pre-exec" and "post-exec" hooks.

 - Look at the output of `zenlog basic-bash-setup` and figure it out.

(However if your shell's command line syntax is far from Posix shell's, then Zenlog may not be able to extract command names 
property and you may not get "per command" output links.)

## Log file structure

By default, log files are stored in `$HOME/zenlog/`, with the following structure:

```
 +--SAN/ # "Sanitized" log -- outout with ESC sequences removed, for easy grepping.
 |  +--YEAR
 |     +--MONTH
 |        +--DAY
 |           +--log files...
 | 
 +--RAW/ # "Raw", or the original output
 |  +--YEAR... (same structure)
 |
 +--ENV/ # "Env", or metadata.
 |  +--YEAR... (same structure)
 |
 +--cmds/ # Per command output.
 |  +--cat
 |  |  +--SAN
 |  |  |  +--YEAR... (same structure)
 |  |  +--RAW
 |  |  +--ENV
 |  |  |--S  # This link contains the last sanitized output from cat(1).
 |  |  |--SS
 |  |  |--...
 |  |  |--R
 |  |  |--RR
 |  |  |--...
 |  |  |--E
 |  |  |--EE
 |  |  |--...
 |  |  
 |  +--ls
 |  :
 | 
 +--pids/  # Per-pid, or "session", output.
 |  + (same strucute)
 | 
 +--tags/  # Per "tag" output.
 |  + (same strucute)
 | 
 |--S  # Link to the last command output.
 |--SS # Link to the second last command output.
 |--...
 |--R
 |--RR
 |--...
 |--E
 |--EE
 |--...
 
```

 - "RAW" log files contain the original output, including all the escape sequences. It's authentic
   but hard to grep.
    
 - "SAW" log files contain the original output with escape sequences stripped out, so easy to grep.
   (Note Zenlog only recognizes often-used escape sequences. Uncommon escape sequences may
   still be left.)
 
 - "ENV" log files contain various meta inforamtion such as the current directory, execution time,
   git branch, etc.
   
 - "S" is a symbolic link to the most recent SAN log file. "R" for RAW, "E" for ENV.
 - "SS", "RR", "EE" are links to the second most log files.
 
 - Zenlog also creates symbolic links for each command and "sessions".
   For example `"$ZENLOG_DIR/pids/$ZENLOG_PID/S"` is a link to the most recent SAN log file
   *on the current shell*. Conversely, `"$ZENLOG_DIR/S"` is the most recent command, which may
   be from a different shell.
   
### Log "tagging"

If you run a command with a comment, for example:
```bash
$ make -B # full build 
```
then Zenlog creates symbolic links in the `tags/` directory too, so `$ZENLOG_DIR/tags/full_build/S`
will be a symbolic link to the most recent "full build" output.

## Advanced customization

 - If you do not want to log output of a specific command (e.g. it doesn't really make sense
   to keep all output from `vi`, `emacs`), you can specify it in
   [`~/.zenlog.toml`](dot_zenlog.toml).
   
   - By default, output from any `zenlog` subcommands will *not* be saved.      

## Useful subcommands

 - `zenlog purge-log [-p DAYS] [-y] [-P]`
   - Removes all log files older than `DAYS` days.
     
     `-y` to execute without a [y/n] prompt.
     
     `-P` for dry-run.    

 - `zenlog du [du(1) options]`
   - Run `du(1)` over the log directory.


 - `zenlog history [-e] [-r] [-n Nth] [-p PID]`
   - Print recent log file names.
   
     `-e` to show the `ENV` log file name instead of `SAN`.
     
     `-r` to show the `RAW` log file name instead of `SAN`.
     
     `-n Nth` Show Nth most recent log file name.
     
       - Note: When you're using this command from a script, the previous command output
        is `-n 1`. But if you're using `zenlog history` from a command that's bound to a hot key
        on the command line, `-n 0` refers to the the previous output. 
     
 - See [this directory](subcommands/) for more (external) subcommands.
   [This file](zenlog/builtins/builtins.go) contains more "buildin" subcommands.  

[See also the readme of the old version.](https://github.com/omakoto/zenlog)

