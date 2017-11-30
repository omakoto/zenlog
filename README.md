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

## Log directory structure

For each command, Zenlog creates two two log files and an optional `ENV` log file. One is called "RAW", which is the original output. However because command output often contains colors, RAW log files are hard to read on an editor and also search. So Zenlog also creates anther log file "SAN" (SANitized), which has most standard ANSI sequences removed.

SAN log files have names like the following. Basically log filenames contain the timestamp, the shortened command line, as well as the comment part in the command line if exists.
```
/home/USER/zenlog/SAN/2017/10/14/14-25-10.778-15193_+ls_-l.log
```

RAW log files have the exact same filename, except they're stored in the `RAW` directory.
```
/home/USER/zenlog/RAW/2017/10/14/14-25-10.778-15193_+ls_-l.log
```

Optionally, if `start-command-with-env` is used instead of `start-command` when starting a command, Zenlog will create `ENV` log files which captures meta information such as the command start time, finish time, status code, environmental variables, etc.

To allow easier access to log files, Zenlog creates a lot of symlinks. (The idea came from [Tanlog](https://github.com/shinh/test/blob/master/tanlog.rb))

```
/home/USER/zenlog/P # The last command SAN output.
/home/USER/zenlog/R # The last command RAW output.
/home/USER/zenlog/PP # Second last SAN log
/home/USER/zenlog/PPP # Third last SAN log
  :
/home/USER/zenlog/RR # Second last RAW log
/home/USER/zenlog/RRR # Third last RAW log
  :
/home/USER/zenlog/cmds/ls/P # The last SAN output from ls
/home/USER/zenlog/cmds/ls/R # The last RAW output from ls
  :
/home/USER/zenlog/tags/TAGNAME/P # The last SAN output with TAGNAME
/home/USER/zenlog/tags/TAGNAME/R # The last RAW output with TAGNAME
  :
```

`TAGNAME` is a comment in command line.

So, for example, if you run the following command:
```
$ cat /etc/fstab | sed '/^#/d' # fstab with comments removed
```

You'll get the regular SAN/RAW log files, as well as the following symlinks:
```
/home/USER/zenlog/cmds/cat/P
/home/USER/zenlog/cmds/cat/R

/home/USER/zenlog/cmds/sed/P
/home/USER/zenlog/cmds/sed/R

/home/USER/zenlog/tags/fstab_comments_removed/P
/home/USER/zenlog/tags/fstab_comments_removed/R
```

`zenlog history` shows the most recent SAN log filenames on the current shell (add `-r` to get RAW names):

```
$ zenlog history
/home/USER/zenlog/SAN/2017/10/14/14-26-54.320-15193_+fstab_comments_removed_+cat_etc_fstab_sed_^#_d_#_fstab_c.log
/home/USER/zenlog/SAN/2017/10/14/14-27-26.700-15193_+zenlog_history.log
```

## Subcommands

### Directive commands

* `zenlog start-command COMMANDLINE`

    Run this command in the pre-exec hook to have Zenlog start logging. COMMANDLINE can be obtained with `bash_last_command` which comes with `zenlog she-helper`, which is explained below.

* `zenlog start-command-with-env "$(bash_dump_env)" COMMANDLINE`

    If this command is used instead of `start-command`, Zenlog will also create `ENV` log files which capture extra information such as command start time, finish time and shell/environmental variables. The first parameter can be the output of
    the `bash_dump_env` command, which dumps the git current branch and the environmental variables, which is a shell function that's in `zenlog sh-helper`, but really any string can be passed as the first argument, and it'll be logged.

* `zenlog stop-log [-n] [exit-status]`

    Run this command in the post-exec hook to have Zenlog stop logging.

    It is guaranteed that when this command returns, both SAN and RAW log files have been all written and closed. So, for example, counting the number of lines with `wc -l` is safe.

    * If `-n` is given, this will print the number of lines in the last log file.
    * If exit-status is given, it'll be written in the `ENV` log file.


### History commands

* `zenlog history [-r]`  Print most recent log filenames.

Example:
```
$ zenlog history
/zenlog/SAN/2017/10/14/16-06-20.773-01908_+ls_etc_.log
/zenlog/SAN/2017/10/14/16-06-32.075-01908_+cat_etc_passwd.log
/zenlog/SAN/2017/10/14/16-06-40.080-01908_+zenlog_history.log

$ zenlog history -r
/zenlog/RAW/2017/10/14/16-06-20.773-01908_+ls_etc_.log
/zenlog/RAW/2017/10/14/16-06-32.075-01908_+cat_etc_passwd.log
/zenlog/RAW/2017/10/14/16-06-40.080-01908_+zenlog_history.log
/zenlog/RAW/2017/10/14/16-07-02.976-01908_+zenlog_history_-r.log
```

* `zenlog history [-r] -n NTH`  Print NTH last log filename.

```
$ cat /etc/passwd
  :
$ zenlog history -n 1
/zenlog/SAN/2017/10/14/16-08-32.236-01908_+cat_etc_passwd.log
```

* `zenlog open-last-log [-r]` Open last log file. `-r` to open RAW instead of SAN.
* `zenlog open-current-log [-r]` Open current log file. `-r` to open RAW instead of SAN.

* `zenlog last-log [-r]` Print last log file name. `-r` to show RAW name instead of SAN.
    - See also the `zenlog_last_log` shell function, which is an order of magnitude faster.

* `zenlog current-log [-r]` Print current log file name. `-r` to show RAW name instead of SAN.
    - See also the `zenlog_current_log` shell function, which is an order of magnitude faster.

    `last-log` and `current-log` are useful for scripting. `current-log` is useful when executing a command *on* the command line prompt. For example, on Bash, you can define a hotkey to launch a command with `bind -x`. Using this, `bind -x '"\e1": "zenlog open-current-log"'` allows you to open the last log file with pressing `ALT-1`.


### Shell helper functions

The following are *shell functions (not Zenlog subcommands)* that can be installed with the following command:
```
. <(zenlog sh-helper)
```

The included functions are:

* `184` executes a command passed as an argument without logging.

    Example: This will print the content of the `a_huge_file_that_not_worth_logging` file without logging it.
```
$ 184 cat a_huge_file_that_not_worth_logging
```

* `186` executes a command with logging, even if a command contains NO_LOG commands. ("man" is included in the default no log list, so normally the output won't be logged. See `ZENLOG_ALWAYS_NO_LOG_COMMANDS` below.)

    Example: This runs `man bash` with logging it.
```
$ 186 man bash
```

* `bash_last_command` shows the most recent command line. Intended to be used with `start-command`. See [the sample bash config file](shell/zenlog.bash).

* `zenlog_in_zenlog`, `zenlog_last_log` and `zenlog_current_log` provide the same functionalities as `zenlog in-zenlog`, `zenlog last-log` and `zenlog current-log`, except they are a lot faster since they're shell functions. Useful when using for the prompt.

### Scripting helper commands

* `zenlog in_zenlog`

    Return success status if in a zenlog session.

    - See also the `zenlog_in_zenlog` shell function, which is an order of magnitude faster.

* `zenlog fail_if_in_zenlog`

    Return success status only if in a zenlog session; otherwise it prints an error message.

* `zenlog fail_unless_in_zenlog`

    Return success status only if *not* in a zenlog session; otherwise it prints an error message.

* `zenlog outer_tty`

    Print the external TTY name, which can be used to print something in the terminal without logging it.

    *Note* if you write to the TTY directly, you'll have to explicitly use CRLF instead of LF. Note the `-e` option and the end of line `\r` in the following example.

    Example:
```
$ echo -e "This is shown but not logged\r" > $(zenlog outer_tty)
```

* `zenlog logger_pipe`

    Print the pipe name to the logger, which can be used to log something without printing it.

    Example:
```
$ echo "This is logged but not shown" > $(zenlog logger_pipe)
```

* `zenlog write_to_outer`

    Eat from STDIN and write to the TTY directly.

    Example:
```
$ echo "This is shown but not logged" | zenlog write_to_outer
```


* `zenlog write_to_logger`

    Eat from STDIN and write to the logger directly. Note the bellow example actually doesn't really work because `zenlog` is a no log command and the log file content will be omitted. This is normally used in a script.

    Example:
```
$ echo "This is logged but not shown" | zenlog write_to_logger
```

### Other commands

* `zenlog purge-log -p DAYS` Remove logs older than DAYS days.

* `zenlog du [du options]` Run du(1) on the log directory.

* `zenlog free-space` Show the free space size of the disk that contains the log directory in bytes.


## Configuration

### Environmental variables

* `ZENLOG_DIR`

    Log directory. Default is `$HOME/zenlog`.

* `ZENLOG_PREFIX_COMMANDS`

    A regex that matches command names that are considered "prefix", which will be ignored when Zenlog detects command names. For example, if `time` is a prefix command, when you run `time cc ....`, the log filename will be `"cc"`, not `"time"`.

    The default is `"(?:command|builtin|time|sudo)"`

* `ZENLOG_ALWAYS_NO_LOG_COMMANDS`

    A regex that matches command names that shouldn't be logged. When Zenlog detect these commands, it'll create log files, but the content will be omitted.

    The default is `"(?:vi|vim|man|nano|pico|less|watch|emacs|zenlog.*)"`

### RC file

* `$HOME/.zenlogrc.rb`

    If this file exists, Zenlog loads it before starting a session.

    If you star Zenlog directly form a terminal application, Zenlog starts before the actual login shell starts, so you can't configure it with the shell's RC file. Instead you can configure environmental variables in this file.

## History

* v0 -- the original version. It evolved from a proof-of-concept script, and was bash-perl hybrid, because I wanted to keep everything in a single file, and I also kinda liked the ugliness of the hybrid script. The below script shows the basic structure. ZENLOG_TTY was (and still is) used to check if the current terminal is within Zenlog or not. (We can't just use "export IN_ZENLOG=1" because if a GUI terminal app is launched from within a Zenlog session, it'd inherit IN_ZENLOG. So instead Zenlog stores the actual TTY name and make sure the current TTY name is the same as the stored one. Technically this could still misjudge when the Zenlog session has already been closed and the TTY name is reused, but that's probably too rare to worry about.)

```bash
#!/bin/bash

# initialization, etc in bash

script -qcf 'ZENLOG_TTY=$(tty) exec bash -l' >(perl <(cat <<'EOF'
# Logger script in perl
EOF
) )
```

* v1 -- The hybrid script was getting harder and harder to maintain and was also ugly, so I finally gave up and split it into multiple files. Also subcommands were now extracted to separate files. v1 still had both the Bash part and the Perl part.

 * v2 -- Rewrote entirely in Perl. No more Bash, except in external subcommands.

 * v3 -- First attempt to rewrite in Ruby, but I soon got bored and it didn't happen.

 * v4 -- v2 was still ugly and hard to improve, so finally rewrote in Ruby. This version has a lot better command line parser, for example, which is used to detect command names in a command line. v2's parser was very hacky so it could mis-parse.

# Caveats

* It is not recommended to set zenlog as a login shell, because once something goes wrong, it'd be very hard to recover. Instead, it's recommended to launch it from a terminal program as a custom start command. This way, if Zenlog stops working for whatever reason, you'll be able to switch to a normal shell just by changing the terminal app's setting.

# Relevant work

## A2H ANSI to HTML converter
[A2H](https://github.com/omakoto/a2h-rs) can be used to convert RAW log files into HTML.

## Compromise zenlog completion

[Compromise](https://github.com/omakoto/compromise) has [zenlog shell completion](https://github.com/omakoto/compromise/blob/master/compromise-zenlog.rb) (WIP).
