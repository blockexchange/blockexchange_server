package worldedit

import (
	"fmt"
	"math"
	"strconv"
)

// This file implements a small, non-executing parser for the subset of Lua
// that worldedit / minetest.serialize() emits: a `return` statement followed
// by a table constructor containing tables, strings, numbers, booleans and nil.
//
// No code is ever executed: function calls, operators (except a unary minus on
// number literals), variables, comments and long strings are rejected.

const (
	// maxDepth bounds the table nesting depth (the worldedit format needs 5:
	// root -> entry -> meta -> inventory -> list)
	maxDepth = 16
	// maxTableItems bounds the number of items in a single nested table
	maxTableItems = 1 << 16
	// maxStringLength bounds the length of a single string literal (decoded)
	maxStringLength = 1 << 20
	// maxNumberLength bounds the length of a single number literal
	maxNumberLength = 64
)

// luaTable is the parsed representation of a lua table constructor
type luaTable struct {
	// positional values: {"a", "b"}
	arr []any
	// string keys: {x=1} or {["x"]=1}
	hash map[string]any
	// explicit integer keys: {[1]="a"}
	idx map[int64]any
	// total number of items in the constructor
	count int
}

// get returns the value of a string key or nil if not present
func (t *luaTable) get(key string) any {
	if t.hash == nil {
		return nil
	}
	return t.hash[key]
}

// list returns the sequence 1..n of the table, following lua semantics
// (positional values take precedence over explicit integer keys)
func (t *luaTable) list() ([]any, error) {
	if len(t.hash) > 0 {
		return nil, fmt.Errorf("expected a list, found string keys")
	}
	if len(t.idx) == 0 {
		return t.arr, nil
	}

	res := append([]any{}, t.arr...)
	used := 0
	for k := range t.idx {
		if k >= 1 && k <= int64(len(t.arr)) {
			// shadowed by a positional value
			used++
		}
	}
	for {
		v, found := t.idx[int64(len(res)+1)]
		if !found {
			break
		}
		res = append(res, v)
		used++
	}
	if used != len(t.idx) {
		return nil, fmt.Errorf("list contains holes or non-positive indices")
	}
	return res, nil
}

type luaParser struct {
	data  []byte
	pos   int
	depth int
}

func newLuaParser(data []byte, offset int) *luaParser {
	return &luaParser{data: data, pos: offset}
}

// errorf returns an error annotated with the current line and column
func (p *luaParser) errorf(format string, args ...any) error {
	line, col := 1, 1
	end := min(p.pos, len(p.data))
	for i := 0; i < end; i++ {
		if p.data[i] == '\n' {
			line++
			col = 1
		} else {
			col++
		}
	}
	return fmt.Errorf("line %d, col %d: %s", line, col, fmt.Sprintf(format, args...))
}

func (p *luaParser) eof() bool {
	return p.pos >= len(p.data)
}

func (p *luaParser) peek() byte {
	if p.eof() {
		return 0
	}
	return p.data[p.pos]
}

func isSpace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r' || c == '\v' || c == '\f'
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func isHexDigit(c byte) bool {
	return isDigit(c) || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')
}

func isIdentStart(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

func isIdentChar(c byte) bool {
	return isIdentStart(c) || isDigit(c)
}

func (p *luaParser) skipSpace() {
	for !p.eof() && isSpace(p.data[p.pos]) {
		p.pos++
	}
}

func (p *luaParser) expect(c byte) error {
	p.skipSpace()
	if p.peek() != c {
		return p.unexpected(fmt.Sprintf("'%c'", c))
	}
	p.pos++
	return nil
}

func (p *luaParser) unexpected(expected string) error {
	if p.eof() {
		return p.errorf("unexpected end of input, expected %s", expected)
	}
	return p.errorf("unexpected character %q, expected %s", p.peek(), expected)
}

func (p *luaParser) readIdent() string {
	start := p.pos
	for !p.eof() && isIdentChar(p.data[p.pos]) {
		p.pos++
	}
	return string(p.data[start:p.pos])
}

// expectKeyword consumes the given keyword (e.g. "return")
func (p *luaParser) expectKeyword(kw string) error {
	p.skipSpace()
	if !isIdentStart(p.peek()) {
		return p.unexpected("'" + kw + "'")
	}
	start := p.pos
	if ident := p.readIdent(); ident != kw {
		p.pos = start
		return p.errorf("unexpected identifier %q, expected '%s'", ident, kw)
	}
	return nil
}

// expectEOF ensures only whitespace follows
func (p *luaParser) expectEOF() error {
	p.skipSpace()
	if !p.eof() {
		return p.unexpected("end of input")
	}
	return nil
}

// parseValue parses a single value: table, string, number, boolean or nil
func (p *luaParser) parseValue() (any, error) {
	p.skipSpace()
	c := p.peek()
	switch {
	case c == '{':
		return p.parseTable()
	case c == '"' || c == '\'':
		return p.parseString()
	case isDigit(c) || c == '-' || c == '.':
		return p.parseNumber()
	case isIdentStart(c):
		start := p.pos
		switch ident := p.readIdent(); ident {
		case "true":
			return true, nil
		case "false":
			return false, nil
		case "nil":
			return nil, nil
		default:
			p.pos = start
			return nil, p.errorf("unexpected identifier %q", ident)
		}
	default:
		return nil, p.unexpected("a value")
	}
}

// parseTableItems parses a table constructor and calls fn for every item,
// key is nil for positional items
func (p *luaParser) parseTableItems(fn func(key, value any) error) error {
	if err := p.expect('{'); err != nil {
		return err
	}
	p.depth++
	if p.depth > maxDepth {
		return p.errorf("tables nested deeper than %d levels", maxDepth)
	}

	for {
		p.skipSpace()
		if p.peek() == '}' {
			p.pos++
			p.depth--
			return nil
		}

		var key any
		switch c := p.peek(); {
		case c == '[' && p.pos+1 < len(p.data) && (p.data[p.pos+1] == '[' || p.data[p.pos+1] == '='):
			return p.errorf("long strings are not supported")
		case c == '[':
			// [expr] = value
			p.pos++
			k, err := p.parseValue()
			if err != nil {
				return err
			}
			if err = p.expect(']'); err != nil {
				return err
			}
			if err = p.expect('='); err != nil {
				return err
			}
			if k == nil {
				return p.errorf("table key is nil")
			}
			key = k
		case isIdentStart(c):
			// name = value, or a keyword value (true/false/nil)
			start := p.pos
			ident := p.readIdent()
			p.skipSpace()
			if p.peek() == '=' && (p.pos+1 >= len(p.data) || p.data[p.pos+1] != '=') {
				p.pos++
				key = ident
			} else {
				p.pos = start
			}
		}

		value, err := p.parseValue()
		if err != nil {
			return err
		}
		if err = fn(key, value); err != nil {
			return err
		}

		p.skipSpace()
		switch p.peek() {
		case ',', ';':
			p.pos++
		case '}':
			// handled at loop start
		default:
			return p.unexpected("',' or '}'")
		}
	}
}

func (p *luaParser) parseTable() (*luaTable, error) {
	t := &luaTable{}
	err := p.parseTableItems(func(key, value any) error {
		t.count++
		if t.count > maxTableItems {
			return p.errorf("table has more than %d items", maxTableItems)
		}

		switch k := key.(type) {
		case nil:
			t.arr = append(t.arr, value)
		case string:
			if t.hash == nil {
				t.hash = map[string]any{}
			}
			t.hash[k] = value
		case float64:
			if k != math.Trunc(k) || math.Abs(k) > 1<<53 {
				return p.errorf("unsupported non-integer table key %v", k)
			}
			if t.idx == nil {
				t.idx = map[int64]any{}
			}
			t.idx[int64(k)] = value
		default:
			return p.errorf("unsupported table key type %T", key)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (p *luaParser) parseNumber() (float64, error) {
	negative := false
	if p.peek() == '-' {
		negative = true
		p.pos++
		p.skipSpace()
	}

	start := p.pos
	var n float64
	if p.peek() == '0' && p.pos+1 < len(p.data) && (p.data[p.pos+1] == 'x' || p.data[p.pos+1] == 'X') {
		// hex integer
		p.pos += 2
		for !p.eof() && isHexDigit(p.data[p.pos]) {
			p.pos++
		}
		token := string(p.data[start+2 : p.pos])
		if token == "" || len(token) > 16 {
			return 0, p.errorf("invalid hex number")
		}
		u, err := strconv.ParseUint(token, 16, 64)
		if err != nil {
			return 0, p.errorf("invalid hex number: %v", err)
		}
		n = float64(u)
	} else {
		digits := 0
		for !p.eof() && isDigit(p.data[p.pos]) {
			p.pos++
			digits++
		}
		if p.peek() == '.' {
			p.pos++
			for !p.eof() && isDigit(p.data[p.pos]) {
				p.pos++
				digits++
			}
		}
		if digits == 0 {
			return 0, p.errorf("invalid number")
		}
		if c := p.peek(); c == 'e' || c == 'E' {
			p.pos++
			if c := p.peek(); c == '+' || c == '-' {
				p.pos++
			}
			expDigits := 0
			for !p.eof() && isDigit(p.data[p.pos]) {
				p.pos++
				expDigits++
			}
			if expDigits == 0 {
				return 0, p.errorf("invalid number exponent")
			}
		}
		if p.pos-start > maxNumberLength {
			return 0, p.errorf("number literal too long")
		}
		f, err := strconv.ParseFloat(string(p.data[start:p.pos]), 64)
		if err != nil {
			return 0, p.errorf("invalid number: %v", err)
		}
		n = f
	}

	if !p.eof() && isIdentChar(p.data[p.pos]) {
		return 0, p.errorf("malformed number")
	}
	if negative {
		n = -n
	}
	return n, nil
}

func (p *luaParser) parseString() (string, error) {
	quote := p.data[p.pos]
	p.pos++

	// fast path: no escape sequences
	start := p.pos
	for i := p.pos; i < len(p.data); i++ {
		c := p.data[i]
		if c == quote {
			if i-start > maxStringLength {
				return "", p.errorf("string longer than %d bytes", maxStringLength)
			}
			p.pos = i + 1
			return string(p.data[start:i]), nil
		}
		if c == '\\' || c == '\n' || c == '\r' {
			break
		}
	}

	buf := make([]byte, 0, 64)
	for {
		if len(buf) > maxStringLength {
			return "", p.errorf("string longer than %d bytes", maxStringLength)
		}
		if p.eof() {
			return "", p.errorf("unfinished string")
		}
		c := p.data[p.pos]
		switch {
		case c == quote:
			p.pos++
			return string(buf), nil
		case c == '\n' || c == '\r':
			return "", p.errorf("unfinished string")
		case c != '\\':
			buf = append(buf, c)
			p.pos++
			continue
		}

		// escape sequence
		p.pos++
		if p.eof() {
			return "", p.errorf("unfinished string")
		}
		e := p.data[p.pos]
		p.pos++
		switch e {
		case 'a':
			buf = append(buf, '\a')
		case 'b':
			buf = append(buf, '\b')
		case 'f':
			buf = append(buf, '\f')
		case 'n':
			buf = append(buf, '\n')
		case 'r':
			buf = append(buf, '\r')
		case 't':
			buf = append(buf, '\t')
		case 'v':
			buf = append(buf, '\v')
		case '\\', '"', '\'':
			buf = append(buf, e)
		case '\n', '\r':
			// escaped line break (emitted by string.format("%q")),
			// "\r\n" and "\n\r" count as one line break
			buf = append(buf, '\n')
			if n := p.peek(); (n == '\n' || n == '\r') && n != e {
				p.pos++
			}
		case 'x':
			if p.pos+2 > len(p.data) || !isHexDigit(p.data[p.pos]) || !isHexDigit(p.data[p.pos+1]) {
				return "", p.errorf("invalid hex escape sequence")
			}
			v, _ := strconv.ParseUint(string(p.data[p.pos:p.pos+2]), 16, 8)
			buf = append(buf, byte(v))
			p.pos += 2
		default:
			if !isDigit(e) {
				return "", p.errorf("invalid escape sequence '\\%c'", e)
			}
			// decimal escape: up to 3 digits
			v := int(e - '0')
			for i := 0; i < 2 && !p.eof() && isDigit(p.data[p.pos]); i++ {
				v = v*10 + int(p.data[p.pos]-'0')
				p.pos++
			}
			if v > 255 {
				return "", p.errorf("decimal escape too large")
			}
			buf = append(buf, byte(v))
		}
	}
}
