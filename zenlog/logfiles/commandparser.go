package logfiles

// Parse a command line and extract executable names and a comment out of it.

import (
	"github.com/omakoto/zenlog-go/zenlog/config"
	"github.com/omakoto/zenlog-go/zenlog/shell"
	"github.com/omakoto/zenlog-go/zenlog/util"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	noLogPrefix    = "184"
	forceLogPrefix = "186"

	whitespaceChars = " \t\n\r\v"
)

var (
	reWordSplitter = util.NewLazyRegexp(`\s+`)
)

// Command represents a single running command executed by the user on the shell.
type Command struct {
	// Trim()'ed version of the command line.
	CommandLine string

	// Command names in a command line. e.g. ["cat", "grep"]
	ExeNames []string

	// Comment in a command line.
	Comment string

	// Whether the command line contained a no log command or had 184.
	// When 186 is detected, the output will always be logged.
	NoLog bool
}

const (
	defaultCommandSplitter = `\s*(?:\;|\|\&|\|\||\||&&|\(|\)|{|})\s*`
	defaultCommentSplitter = `(?:^|[\;\|\&\(\)\{\}]|\s+)\#\s*`
)

// Take a command line and split into the command part and the comment part.
// e.g. "this is command # comment"
// Note when # appears in the middle of a word, it won't start a comment.
// i.g. "this is comma#nd" is a 3 word command and the # mark doesn't start a comment.
func splitComment(config *config.Config, commandLine string) (string, string) {
	pat := util.FirstNonEmpty(config.CommentSplitter, defaultCommentSplitter)
	re := regexp.MustCompile(pat)
	vals := re.Split(commandLine, 2)
	if len(vals) == 2 {
		return vals[0], vals[1]
	}
	return vals[0], ""
}

func splitCommands(config *config.Config, commandLine string) []string {
	re := regexp.MustCompile(util.FirstNonEmpty(config.CommandSplitter, defaultCommandSplitter))
	return re.Split(commandLine, -1)
}

func extractCommandsWithRegex(config *config.Config, commandLine string) (commands [][]string, comment string) {

	pipeLine, comment := splitComment(config, commandLine)

	commands = make([][]string, 0, 16)

	for _, command := range splitCommands(config, pipeLine) {
		commands = append(commands, reWordSplitter.Pattern().Split(command, -1))
	}
	return
}

func extractCommandsWithParser(config *config.Config, commandLine string) (commands [][]string, comment string) {
	commands = make([][]string, 0, 16)
	comment = ""

	tokens := shell.Split(commandLine)
	if len(tokens) == 0 {
		return
	}
	last := tokens[len(tokens)-1]
	if strings.HasPrefix(last, "#") {
		comment = last[1:]
		comment = strings.Trim(comment, " \t\r\n")
		tokens = tokens[0 : len(tokens)-1]
	}

	current := make([]string, 0, 16)

	push := func() {
		if len(current) > 0 {
			commands = append(commands, current)
			current = make([]string, 0, 16)
		}
	}

	for _, word := range tokens {
		if shell.IsCommandSeparator(word) {
			push()
			continue
		}
		current = append(current, word)
	}
	push()
	return
}

func extractCommands(config *config.Config, commandLine string) (commands [][]string, comment string) {
	if !config.UseExperimentalCommandParser {
		return extractCommandsWithRegex(config, commandLine)
	}
	return extractCommandsWithParser(config, commandLine)
}

// ParseCommandLine takes a command line, and extracts a list of the commands and the comment out of it.
// e.g. "/bin/cat /etc/passwd | grep xxx # find xxx" -> ["cat", "grep"] "find xxx"
func ParseCommandLine(config *config.Config, commandLine string) *Command {
	// Save command.
	ret := Command{}
	ret.CommandLine = strings.Trim(commandLine, whitespaceChars)

	// Tokenize.
	commands, comment := extractCommands(config, commandLine)

	ret.Comment = comment

	// Get command names, and check 184/186.
	prefixCommands := regexp.MustCompile("^" + config.PrefixCommands + "$")
	alwaysNoLogCommands := regexp.MustCompile("^" + config.AlwaysNoLogCommands + "$")

	exes := make([]string, 0, 8)

	noLogDetected := false
	forceLogDetected := false

	for i, command := range commands {
		for _, w := range command {
			switch w {
			case noLogPrefix:
				noLogDetected = true
				continue
			case forceLogPrefix:
				forceLogDetected = true
				continue
			}
			if prefixCommands.MatchString(w) {
				continue
			}

			// Let's only use the first command for auto-184.
			if i == 0 && alwaysNoLogCommands.MatchString(w) {
				noLogDetected = true
			}
			exes = append(exes, filepath.Base(w))
			break
		}
	}

	ret.ExeNames = exes
	ret.NoLog = !forceLogDetected && noLogDetected
	return &ret
}
