package util

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
	"time"
)

const (
	//commandMarker = "!zenlog:"
	//escape         = '\\'

	commandMarker = "\x1b\x01\x09\x07\x03\x02\x05zenlog:"
	escape        = '\x1b'
	separator     = ' '
)

var (
	commandMarkerBytes = []byte(commandMarker)
)

func _addHexDigit(b *bytes.Buffer, v uint8) {
	if v <= 9 {
		b.WriteRune(rune('0' + v))
	} else if 10 <= v && v <= 15 {
		b.WriteRune(rune('a' + (v - 10)))
	} else {
		panic(fmt.Sprintf("Invalid value %d", v))
	}
}

func _hexToInt(v uint8) uint8 {
	if '0' <= v && v <= '9' {
		return v - '0'
	} else if 'a' <= v && v <= 'f' {
		return v - 'a' + 10
	} else {
		panic(fmt.Sprintf("Invalid value %d", v))
	}
}

func _encodeSingle(b *bytes.Buffer, arg string) {
	for j := 0; j < len(arg); j++ {
		ch := arg[j]
		switch ch {
		case '\n', escape, separator:
			b.WriteByte(escape)
			_addHexDigit(b, (ch>>4)&0xf)
			_addHexDigit(b, (ch>>0)&0xf)
		default:
			b.WriteByte(ch)
		}
	}
}

// Encode a string array into a single command string.
func Encode(args []string) string {
	var b bytes.Buffer

	b.WriteString(commandMarker)

	for i := 0; i < len(args); i++ {
		if i >= 1 {
			b.WriteRune(separator)
		}
		arg := args[i]
		_encodeSingle(&b, arg)
	}
	b.WriteByte('\n')

	return b.String()
}

func _decodeSingle(arg string) string {
	var b bytes.Buffer

	for i := 0; i < len(arg); i++ {
		ch := arg[i]
		switch ch {
		case escape:
			b.WriteByte((_hexToInt(arg[i+1]) << 4) + _hexToInt(arg[i+2]))
			i += 2
		default:
			b.WriteByte(ch)
		}
	}

	return b.String()
}

func _decode(s string) []string {
	if s == "" {
		return make([]string, 0)
	}
	eargs := strings.Split(s, string(separator))
	ret := make([]string, len(eargs))
	for i, a := range eargs {
		ret[i] = _decodeSingle(a)
	}
	return ret
}

// TryDecodeBytes extracts a command string from a given string, and return the prefix part
// (i.e. the part followed by the command string) and the split up command
// arguments.
func TryDecodeBytes(s []byte) (success bool, prefix []byte, args []string) {
	pos := bytes.Index(s, commandMarkerBytes)
	if pos < 0 {
		success = false
		prefix = nil
		args = nil
		return
	}
	success = true
	prefix = s[0:pos]

	// If the last character is \n, chop it.
	end := len(s) - 1
	if s[end] == '\n' {
		end--
	}

	args = _decode(string(s[pos+len(commandMarker) : end+1]))
	return
}

// WriteToFile writes encoded args to a file.
func WriteToFile(filename string, args []string) error {
	file, err := os.OpenFile(filename, os.O_WRONLY, 0)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(Encode(args))
	return err
}

// ReadFromFile read encoded args from a file. The predicate callback can be used to tell whether a
// received set of args is what's expected or not.
func ReadFromFile(filename string, predicate func(vals []string) bool, timeout time.Duration) ([]string, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	type result struct {
		vals []string
		err  error
	}

	ch := make(chan result)
	defer close(ch)

	bi := bufio.NewReader(file)

	go func() {
		for {
			line, err := bi.ReadBytes('\n')
			if line != nil {
				success, _, args := TryDecodeBytes(line)
				Debugf("reply: %v", args)
				if success && predicate != nil && predicate(args) {
					ch <- result{args, nil}
					return
				}
			}
			if err != nil {
				ch <- result{nil, err}
				return
			}
		}
	}()

	select {
	case res := <-ch:
		return res.vals, res.err
	case <-time.After(timeout):
		return nil, fmt.Errorf("timed out reading from '%s'", filename)
	}
}
