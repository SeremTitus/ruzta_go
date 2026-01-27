package tokenizer

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

type Tokenizer struct {
	source           []rune
	_start           int
	startLine        int
	startColumn      int
	errorStack       []*Token
	pendingNewline   bool
	lastToken        *Token
	lastNewline      *Token
	_current         int
	line             int
	column           int
	position         int
	length           int
	tabSize          int
	parenStack       []rune
}

func NewTokenizer(src string) *Tokenizer {
	return &Tokenizer{
		source:      []rune(src),
		line:        1,
		column:      1,
		length:      len([]rune(src)),
		tabSize:     4,
	}
}

func isDigit(p_char rune) bool {
	return p_char >= '0' && p_char <= '9'
}

func isUnicodeIdentifierStart(r rune) bool {
	return r == '_' || unicode.IsLetter(r)
}

func isUnicodeIdentifierContinue(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
}

func isWhitespace(ch rune) bool {
	return ch == ' ' ||
		ch == 0x00A0 ||
		ch == 0x1680 ||
		(ch >= 0x2000 && ch <= 0x200B) ||
		ch == 0x202F ||
		ch == 0x205F ||
		ch == 0x3000 ||
		ch == 0x2028 ||
		ch == 0x2029 ||
		(ch >= 0x0009 && ch <= 0x000D) ||
		ch == 0x0085
}

func (t *Tokenizer) isAtEnd() bool {
	return t._current >= t.length
}

func (t *Tokenizer) peek(offset int) rune {
	idx := t._current + offset
	if idx >= 0 && idx < t.length {
		return t.source[idx]
	}
	return 0
}

func (t *Tokenizer) advance() rune {
	if t.isAtEnd() {
		return 0
	}
	ch := t.source[t._current]
	t._current++
	t.position++
	t.column++
	return ch
}

func (t *Tokenizer) pushParen(char rune) {
	t.parenStack = append(t.parenStack, char)
}

func (t *Tokenizer) popParen(expected rune) bool {
	if len(t.parenStack) == 0 {
		return false
	}
	actual := t.parenStack[len(t.parenStack)-1]
	t.parenStack = t.parenStack[:len(t.parenStack)-1]
	return actual == expected
}

func (t *Tokenizer) makeParenError(paren rune) *Token {
	if len(t.parenStack) == 0 {
		return t.makeError(fmt.Sprintf("Closing \"%c\" doesn't have an opening counterpart.", paren))
	}
	rError := t.makeError(fmt.Sprintf("Closing \"%c\" doesn't match the opening \"%c\".", paren, t.parenStack[len(t.parenStack)-1]))
	t.parenStack = t.parenStack[:len(t.parenStack)-1]
	return rError
}

func (t *Tokenizer) hasError() bool {
	return len(t.errorStack) > 0
}

func (t *Tokenizer) pushError(msg string) {
	rError := t.makeError(msg)
	t.errorStack = append(t.errorStack, rError)
}

func (t *Tokenizer) popError() *Token {
	rError := t.errorStack[len(t.errorStack)-1]
	t.errorStack = t.errorStack[:len(t.errorStack)-1]
	return rError
}

func (t *Tokenizer) makeToken(tokenType TokenType) *Token {
	token := NewToken(tokenType)
	token.StartLine = t.line
	token.EndLine = t.line
	token.StartColumn = t.column
	token.EndColumn = t.column
	token.Source = t.source[t._start:t._current]
	t.lastToken = token
	return token
}

func (t *Tokenizer) makeLiteral(value interface{}) *Token {
	token := t.makeToken(LITERAL)
	token.Literal = value
	return token
}

func (t *Tokenizer) makeIdentifier(name string) *Token {
	token := t.makeToken(IDENTIFIER)
	token.Literal = name
	return token
}

func (t *Tokenizer) makeError(msg string) *Token {
	token := t.makeToken(ERROR)
	token.Literal = msg
	return token
}


func (t *Tokenizer) skipWhitespace() {
	for {
		switch t.peek(0) {
		case ' ':
			t.advance()

		case '\t':
			t.advance()
			t.column += t.tabSize - 1

		case '\r':
			t.advance()
			if t.peek(0) != '\n' {
				t.pushError("Stray carriage return character in source code.")
				return
			}

		case '\n':
			t.advance()
			t.newline(true)

		case '#':
			t.skipLineComment()

		case '/':
			if t.peek(1) == '/' {
				t.advance()
				t.advance()
				t.skipLineComment()
			} else if t.peek(1) == '*' {
				t.advance()
				t.advance()
				t.skipBlockComment()
			} else {
				return
			}

		default:
			return
		}
	}
}

func (t *Tokenizer) skipLineComment() {
	for t.peek(0) != '\n' && !t.isAtEnd() {
		t.advance()
	}
	if t.isAtEnd() {
		return
	}
	t.advance()
	t.newline(true)
}

func (t *Tokenizer) skipBlockComment() {
	for {
		if t.isAtEnd() {
			t.pushError("Unterminated block comment.")
			return
		}
		if t.peek(0) == '\r' {
			t.advance()
			if t.peek(0) != '\n' {
				t.pushError("Stray carriage return character in source code.")
				return
			}
			t.advance()
			t.newline(true)
			continue
		}
		if t.peek(0) == '\n' {
			t.advance()
			t.newline(true)
			continue
		}
		if t.peek(0) == '*' && t.peek(1) == '/' {
			t.advance()
			t.advance()
			return
		}
		t.advance()
	}
}



func (t *Tokenizer) newline(make bool) {
	// Don't overwrite a previous newline token.
	if make && !t.pendingNewline {
		lineToken := t.makeToken(NEWLINE)
		lineToken.StartLine = t.line
		lineToken.EndLine = t.line
		lineToken.StartColumn = t.column - 1
		lineToken.EndColumn = t.column
		t.pendingNewline = true
		t.lastToken = lineToken
		t.lastNewline = lineToken
	}

	t.line++
	t.column = 1
}

func (t *Tokenizer) number() *Token {
	start := t._current - 1
	first := t.source[start]

	if first == '0' {
		switch t.peek(0) {
		case 'x', 'X', 'b', 'B', 'o', 'O':
			prefix := t.peek(0)
			t.advance()
			base := 10
			switch prefix {
			case 'x', 'X':
				base = 16
			case 'b', 'B':
				base = 2
			case 'o', 'O':
				base = 8
			}
			digits := 0
			for {
				ch := t.peek(0)
				if ch == '_' || ch == ',' {
					t.advance()
					continue
				}
				if isDigitForBase(ch, base) {
					digits++
					t.advance()
					continue
				}
				break
			}
			if digits == 0 {
				return t.makeError("Expected digits after base prefix.")
			}
			raw := string(t.source[start:t._current])
			clean := removeSeparators(raw)[2:]
			value, err := strconv.ParseInt(clean, base, 64)
			if err != nil {
				return t.makeError("Invalid numeric literal.")
			}
			return t.makeLiteral(value)
		}
	}

	sawDot := first == '.'
	if !sawDot {
		for {
			ch := t.peek(0)
			if isDigit(ch) || ch == '_' || ch == ',' {
				t.advance()
			} else {
				break
			}
		}
	}

	if sawDot {
		for {
			ch := t.peek(0)
			if isDigit(ch) || ch == '_' || ch == ',' {
				t.advance()
			} else {
				break
			}
		}
	}

	if !sawDot && t.peek(0) == '.' && isDigit(t.peek(1)) {
		sawDot = true
		t.advance()
		for {
			ch := t.peek(0)
			if isDigit(ch) || ch == '_' || ch == ',' {
				t.advance()
			} else {
				break
			}
		}
	}

	sawExp := false
	if t.peek(0) == 'e' || t.peek(0) == 'E' {
		sawExp = true
		t.advance()
		if t.peek(0) == '+' || t.peek(0) == '-' {
			t.advance()
		}
		if !isDigit(t.peek(0)) {
			return t.makeError("Expected exponent digits after 'e'.")
		}
		for {
			ch := t.peek(0)
			if isDigit(ch) || ch == '_' || ch == ',' {
				t.advance()
			} else {
				break
			}
		}
	}

	raw := string(t.source[start:t._current])
	clean := removeSeparators(raw)
	if sawDot || sawExp {
		value, err := strconv.ParseFloat(clean, 64)
		if err != nil {
			return t.makeError("Invalid numeric literal.")
		}
		return t.makeLiteral(value)
	}
	value, err := strconv.ParseInt(clean, 10, 64)
	if err != nil {
		return t.makeError("Invalid numeric literal.")
	}
	return t.makeLiteral(value)
}

func isDigitForBase(ch rune, base int) bool {
	switch {
	case base <= 10:
		return ch >= '0' && ch < rune('0'+base)
	case ch >= '0' && ch <= '9':
		return true
	case base == 16 && ((ch >= 'a' && ch <= 'f') || (ch >= 'A' && ch <= 'F')):
		return true
	default:
		return false
	}
}

func removeSeparators(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		if r == '_' || r == ',' {
			continue
		}
		b.WriteRune(r)
	}
	return b.String()
}

func (t *Tokenizer) stringToken() *Token {
	quote := t.peek(-1)
	result := ""
	for {
		c := t.peek(0)
		if t.isAtEnd() {
			return t.makeError("Unterminated string")
		}
		if c == quote {
			t.advance()
			break
		}
		if c == '\\' {
			t.advance()
			c = t.peek(0)
			t.advance()
			switch c {
			case 'n':
				result += "\n"
			case 't':
				result += "\t"
			case '\\':
				result += "\\"
			case '"':
				result += "\""
			default:
				result += string(c)
			}
		} else {
			result += string(c)
			t.advance()
		}
	}
	return t.makeLiteral(result)
}

var keywordsByFirst = map[rune]map[string]TokenType{
	'a': {
		"as":  AS,
		"and": AND,
	},
	'b': {
		"break": BREAK,
	},
	'c': {
		"class":    CLASS,
		"const":    CONST,
		"continue": CONTINUE,
	},
	'e': {
		"elif":    ELIF,
		"else":    ELSE,
		"enum":    ENUM,
		"extends": EXTENDS,
	},
	'f': {
		"for": FOR,
		"fn":  FUNCTION,
	},
	'i': {
		"if":     IF,
		"import": IMPORT,
		"in":     IN,
		"is":     IS,
	},
	'm': {
		"match": MATCH,
		"mod":   MOD,
	},
	'n': {
		"not": NOT,
	},
	'o': {
		"or": OR,
	},
	'p': {
		"pass": PASS,
	},
	'r': {
		"return": RETURN,
	},
	's': {
		"self":   SELF,
		"signal": SIGNAL,
	},
	't': {
		"trait": TRAIT,
		"type":  TYPE,
	},
	'u': {
		"uses": USES,
	},
	'v': {
		"var":  VAR,
		"void": VOID,
	},
	'w': {
		"while": WHILE,
		"when":  WHEN,
	},
	'I': {
		"INF": CONST_INF,
	},
	'N': {
		"NAN": CONST_NAN,
	},
	'P': {
		"PI": CONST_PI,
	},
	'T': {
		"TAU": CONST_TAU,
	},
}

const (
	MinKeywordLength = 2
	MaxKeywordLength = 10
)

func (t *Tokenizer) potentialIdentifier() *Token {
	start := t._current - 1
	for isUnicodeIdentifierContinue(t.peek(0)) {
		t.advance()
	}
	name := string(t.source[start:t._current])
	length := len(name)

	if length >= MinKeywordLength && length <= MaxKeywordLength {
		switch length {
		case 4:
			if name == "true" {
				return t.makeLiteral(true)
			}
			if name == "null" {
				return t.makeLiteral(nil)
			}
		case 5:
			if name == "false" {
				return t.makeLiteral(false)
			}
		}
		first := []rune(name)[0]
		if group, ok := keywordsByFirst[first]; ok {
			if tokenType, ok := group[name]; ok {
				newToken := t.makeToken(tokenType)
				newToken.Literal = name
				return newToken
			}
		}
		
	}
	return t.makeIdentifier(name)
}

func (t *Tokenizer) checkVCSMarker(test rune, doubleType TokenType) *Token {
	chars := 2             // two already matched
	next := t.position + 1 // lookahead index

	// Count consecutive matching runes WITHOUT consuming
	for next < t.length && t.source[next] == test {
		chars++
		next++
	}

	if chars >= 7 {
		// VCS conflict marker (<<<<<<<, =======, >>>>>>>)
		for chars > 1 {
			t.advance() // first char already consumed by Scan()
			chars--
		}
		return t.makeToken(VCS_CONFLICT_MARKER)
	}

	// Regular double-character token (==, <<, >>, etc.)
	t.advance() // consume second character
	return t.makeToken(doubleType)
}

func (t *Tokenizer) annotation() *Token {
	if isUnicodeIdentifierStart(t.peek(0)) {
		t.advance()
	} else {
		panic("Expected annotation identifier after \"@\".")
	}
	for isUnicodeIdentifierContinue(t.peek(0)) {
		t.advance()
	}
	annotationToken := t.makeToken(ANNOTATION)
	annotationToken.Literal = string(annotationToken.Source)
	return annotationToken
}

func (t *Tokenizer) Scan() *Token {
	if t.hasError() {
		return t.popError()
	}

	t.skipWhitespace()

	if t.pendingNewline {
		t.pendingNewline = false
		return t.lastNewline
	}

	if t.hasError() {
		return t.popError()
	}

	t._start = t._current
	t.startLine = t.line
	t.startColumn = t.column

	if t.isAtEnd() {
		return t.makeToken(EOF)
	}

	c := t.advance()

	if isDigit(c) {
		return t.number()
	} else if c == 'r' && (t.peek(0) == '"' || t.peek(0) == '\'') {
		// Raw string literals.
		return t.stringToken()
	} else if isUnicodeIdentifierStart(c) {
		return t.potentialIdentifier()
	}

	// Single-char tokens
	switch c {
	case '"', '\'':
		return t.stringToken()
	case '@':
		return t.annotation()
	case '~':
		return t.makeToken(TILDE)
	case ',':
		return t.makeToken(COMMA)
	case ':':
		return t.makeToken(COLON)
	case ';':
		return t.makeToken(SEMICOLON)
	case '$':
		return t.makeToken(DOLLAR)
	case '?':
		return t.makeToken(QUESTION_MARK)
	case '`':
		return t.makeToken(BACKTICK)
	case '(':
		t.pushParen('(')
		return t.makeToken(PARENTHESIS_OPEN)
	case '[':
		t.pushParen('[')
		return t.makeToken(BRACKET_OPEN)
	case '{':
		t.pushParen('{')
		return t.makeToken(BRACE_OPEN)
	case ')':
		if !t.popParen('(') {
			return t.makeParenError(c)
		}
		return t.makeToken(PARENTHESIS_CLOSE)
	case ']':
		if !t.popParen('[') {
			return t.makeParenError(c)
		}
		return t.makeToken(BRACKET_CLOSE)
	case '}':
		if !t.popParen('{') {
			return t.makeParenError(c)
		}
		return t.makeToken(BRACE_CLOSE)
	// Double characters.
	case '!':
		if t.peek(0) == '=' {
			t.advance()
			return t.makeToken(BANG_EQUAL)
		} else {
			return t.makeToken(BANG)
		}
	case '.':
		if t.peek(0) == '.' {
			t.advance()
			if t.peek(0) == '.' {
				t.advance()
				return t.makeToken(PERIOD_PERIOD_PERIOD)
			}
			return t.makeToken(PERIOD_PERIOD)
		} else if isDigit(t.peek(0)) {
			// Number starting with '.'.
			return t.number()
		} else {
			return t.makeToken(PERIOD)
		}
	case '+':
		if t.peek(0) == '=' {
			t.advance()
			return t.makeToken(PLUS_EQUAL)
		} else if isDigit(t.peek(0)) && !t.lastToken.CanPrecedeBinOP() {
			// Number starting with '+'.
			return t.number()
		} else {
			return t.makeToken(PLUS)
		}
	case '-':
		if t.peek(0) == '=' {
			t.advance()
			return t.makeToken(MINUS_EQUAL)
		} else if isDigit(t.peek(0)) && !t.lastToken.CanPrecedeBinOP() {
			// Number starting with '-'.
			return t.number()
		} else if t.peek(0) == '>' {
			t.advance()
			return t.makeToken(FORWARD_ARROW)
		} else {
			return t.makeToken(MINUS)
		}
	case '*':
		if t.peek(0) == '=' {
			t.advance()
			return t.makeToken(STAR_EQUAL)
		} else if t.peek(0) == '*' {
			if t.peek(1) == '=' {
				t.advance()
				t.advance() // Advance both '*' and '='
				return t.makeToken(STAR_STAR_EQUAL)
			}
			t.advance()
			return t.makeToken(STAR_STAR)
		} else {
			return t.makeToken(STAR)
		}
	case '/':
		if t.peek(0) == '=' {
			t.advance()
			return t.makeToken(SLASH_EQUAL)
		} else {
			return t.makeToken(SLASH)
		}
	case '%':
		if t.peek(0) == '=' {
			t.advance()
			return t.makeToken(PERCENT_EQUAL)
		} else {
			return t.makeToken(PERCENT)
		}
	case '^':
		if t.peek(0) == '=' {
			t.advance()
			return t.makeToken(CARET_EQUAL)
		} else if t.peek(0) == '"' || t.peek(0) == '\'' {
			// Node path
			return t.stringToken()
		} else {
			return t.makeToken(CARET)
		}
	case '&':
		if t.peek(0) == '&' {
			t.advance()
			return t.makeToken(AMPERSAND_AMPERSAND)
		} else if t.peek(0) == '=' {
			t.advance()
			return t.makeToken(AMPERSAND_EQUAL)
		} else if t.peek(0) == '"' || t.peek(0) == '\'' {
			// String Name
			return t.stringToken()
		} else {
			return t.makeToken(AMPERSAND)
		}
	case '|':
		if t.peek(0) == '|' {
			t.advance()
			return t.makeToken(PIPE_PIPE)
		} else if t.peek(0) == '=' {
			t.advance()
			return t.makeToken(PIPE_EQUAL)
		} else {
			return t.makeToken(PIPE)
		}

	// Potential VCS conflict markers.
	case '=':
		if t.peek(0) == '=' {
			return t.checkVCSMarker('=', EQUAL_EQUAL)
		} else {
			return t.makeToken(EQUAL)
		}
	case '<':
		if t.peek(0) == '=' {
			t.advance()
			return t.makeToken(LESS_EQUAL)
		} else if t.peek(0) == '<' {
			if t.peek(1) == '=' {
				t.advance()
				t.advance() // Advance both '<' and '='
				return t.makeToken(LESS_LESS_EQUAL)
			} else {
				return t.checkVCSMarker('<', LESS_LESS)
			}
		} else {
			return t.makeToken(LESS)
		}
	case '>':
		if t.peek(0) == '=' {
			t.advance()
			return t.makeToken(GREATER_EQUAL)
		} else if t.peek(0) == '>' {
			if t.peek(1) == '=' {
				t.advance()
				t.advance() // Advance both '>' and '='
				return t.makeToken(GREATER_GREATER_EQUAL)
			} else {
				return t.checkVCSMarker('>', GREATER_GREATER)
			}
		} else {
			return t.makeToken(GREATER)
		}
	default:
		if isWhitespace(c) {
			t.skipWhitespace()
			return t.Scan()
		} else {
			return t.makeError(fmt.Sprintf(`Invalid character "%c"`, c))
		}
	}
}
