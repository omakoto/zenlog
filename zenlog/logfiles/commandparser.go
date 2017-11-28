package logfiles

// Parse a command line and extract executable names and a comment out of it.

import (
	"github.com/omakoto/zenlog-go/zenlog/config"
	"github.com/omakoto/zenlog-go/zenlog/util"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	WHITESPACE = " \t\n\r\v"
)

var (
	reWordSplitter = util.NewLazyRegexp(`\s+`)
)

type Command struct {
	// Trim()'ed version of the command line.
	CommandLine string

	// Command names in a command line. e.g. ["cat", "grep"]
	ExeNames []string

	// Comment in a command line.
	Comment string
}

const (
	DEFAULT_COMMAND_SPLITTER = `\s*(?:\;|\|\&|\|\||\||&&|\(|\)|{|})\s*`
	DEFAULT_COMMENT_SPLITTER = `(?:^|[\;\|\&\(\)\{\}]|\s+)\#\s*`
)

// Take a command line and split into the command part and the comment part.
// e.g. "this is command # comment"
// Note when # appears in the middle of a word, it won't start a comment.
// i.g. "this is comma#nd" is a 3 word command and the # mark doesn't start a comment.
func splitComment(config *config.Config, commandLine string) (string, string) {
	pat := util.FirstNonEmpty(config.CommentSplitter, DEFAULT_COMMENT_SPLITTER)
	re := regexp.MustCompile(pat)
	vals := re.Split(commandLine, 2)
	if len(vals) == 2 {
		return vals[0], vals[1]
	}
	return vals[0], ""
}

func splitCommands(config *config.Config, commandLine string) []string {
	re := regexp.MustCompile(util.FirstNonEmpty(config.CommandSplitter, DEFAULT_COMMAND_SPLITTER))
	return re.Split(commandLine, -1)
}

// Take a command line, and extract a list of the commands and the comment.
// e.g. "/bin/cat /etc/passwd | grep xxx # find xxx" -> ["cat", "grep"] "find xxx"
// TODO Make it actually understand quotes, etc.
func ParseCommandLine(config *config.Config, commandLine string) *Command {
	ret := Command{}
	ret.CommandLine = strings.Trim(commandLine, WHITESPACE)

	pipeLine, comment := splitComment(config, commandLine)
	ret.Comment = comment

	commands := splitCommands(config, pipeLine)

	exes := make([]string, 0, 8)

	for _, command := range commands {
		words := reWordSplitter.Pattern().Split(command, -1)
		if len(words) == 0 || words[0] == "" {
			continue
		}
		exes = append(exes, filepath.Base(words[0]))
	}

	ret.ExeNames = exes
	return &ret
}

//var (
//	SHELL_SPECIAL_CHARS    = []byte(";|&<>(){}\n")
//	SHELL_WHITESPACE_CHARS = []byte(" \t")
//
//	SPECIAL_TOKEN    = SHELL_SPECIAL_CHARS[0:1]
//	WHITESPACE_TOKEN = SHELL_WHITESPACE_CHARS[0:1]
//)
//
//func isSpace(b byte) bool {
//	switch b {
//	case ' ', '\t':
//		return true
//	}
//	return false
//}
//
//func consecutive(data []byte, chars []byte) (found bool, token []byte, nextPos int) {
//	nextPos = 0
//	for {
//		if bytes.Index(chars, data[nextPos:nextPos+1]) >= 0 {
//			break
//		}
//		nextPos++
//	}
//	if nextPos == 0 {
//		found = false
//	} else {
//		found = true
//		token = data[0:nextPos]
//	}
//	return
//}
//
//func ScanShellToken(data []byte, atEOF bool) (advance int, token []byte, err error) {
//	if len(data) == 0 && atEOF {
//		return 0, nil, nil
//	}
//	if found, _, next := consecutive(data, SHELL_SPECIAL_CHARS); found {
//		return next, SPECIAL_TOKEN, nil
//	}
//	if found, _, next := consecutive(data, SHELL_WHITESPACE_CHARS); found {
//		return next, WHITESPACE_TOKEN, nil
//	}
//
//	pos := 0
//	for {
//		if pos >= len(data) || isSpace(data[pos]) {
//			break
//		}
//		pos++
//
//	}
//	return pos, data[0:pos], nil
//}
//
//func ParseCommandLine(commandLine string) Command {
//	//exes := make([]string, 0, 8)
//	scanner := bufio.NewScanner(strings.NewReader(commandLine))
//	scanner.Split(ScanShellToken)
//
//	scanner.Scan()
//
//}
