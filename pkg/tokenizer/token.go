package tokenizer

type CursorPlace int

const (
	CURSOR_NONE CursorPlace = iota
	CURSOR_BEGINNING
	CURSOR_MIDDLE
	CURSOR_END
)

type TokenType int

const (
	EMPTY TokenType = iota
	// Basic
	ANNOTATION
	IDENTIFIER
	LITERAL
	// Comparison
	LESS
	LESS_EQUAL
	GREATER
	GREATER_EQUAL
	EQUAL_EQUAL
	BANG_EQUAL
	// Logical
	AND
	OR
	NOT
	AMPERSAND_AMPERSAND
	PIPE_PIPE
	BANG
	// Bitwise
	AMPERSAND
	PIPE
	TILDE
	CARET
	LESS_LESS
	GREATER_GREATER
	// Math
	PLUS
	MINUS
	STAR
	STAR_STAR
	SLASH
	PERCENT
	// Assignment
	EQUAL
	PLUS_EQUAL
	MINUS_EQUAL
	STAR_EQUAL
	STAR_STAR_EQUAL
	SLASH_EQUAL
	PERCENT_EQUAL
	LESS_LESS_EQUAL
	GREATER_GREATER_EQUAL
	AMPERSAND_EQUAL
	PIPE_EQUAL
	CARET_EQUAL
	// Control flow
	IF
	ELIF
	ELSE
	FOR
	WHILE
	BREAK
	CONTINUE
	PASS
	RETURN
	MATCH
	WHEN
	// Keywords
	AS
	CLASS
	CONST
	ENUM
	EXTENDS
	FUNCTION
	IMPORT
	IN
	IS
	MOD
	SELF
	SIGNAL
	TRAIT
	TYPE
	USES
	VAR
	VOID
	// Punctuation
	BRACKET_OPEN
	BRACKET_CLOSE
	BRACE_OPEN
	BRACE_CLOSE
	PARENTHESIS_OPEN
	PARENTHESIS_CLOSE
	COMMA
	SEMICOLON
	PERIOD
	PERIOD_PERIOD
	PERIOD_PERIOD_PERIOD
	COLON
	DOLLAR
	FORWARD_ARROW
	UNDERSCORE
	// Whitespace
	NEWLINE
	// Constants
	CONST_PI
	CONST_TAU
	CONST_INF
	CONST_NAN
	// Error message improvement
	VCS_CONFLICT_MARKER
	BACKTICK
	QUESTION_MARK
	// Special
	ERROR
	EOF // "EOF" is reserved
	MAX
)

type Token struct {
	Type           TokenType
	Literal        interface{} // Variant in C++, can be string/int/etc.
	StartLine      int
	EndLine        int
	StartColumn    int
	EndColumn      int
	CursorPosition int
	CursorPlace    CursorPlace
	Source         []rune
}

func NewToken(p_type TokenType) *Token {
	return &Token{
		Type:           p_type,
		CursorPosition: -1,
		CursorPlace:    CURSOR_NONE,
	}
}

func (t *Token) GetName() string {
	switch t.Type {
	case EMPTY:
		return "Empty"

	// Basic
	case ANNOTATION:
		return "Annotation"
	case IDENTIFIER:
		return "Identifier"
	case LITERAL:
		return "Literal"

	// Comparison
	case LESS:
		return "<"
	case LESS_EQUAL:
		return "<="
	case GREATER:
		return ">"
	case GREATER_EQUAL:
		return ">="
	case EQUAL_EQUAL:
		return "=="
	case BANG_EQUAL:
		return "!="

	// Logical
	case AND:
		return "and"
	case OR:
		return "or"
	case NOT:
		return "not"
	case AMPERSAND_AMPERSAND:
		return "&&"
	case PIPE_PIPE:
		return "||"
	case BANG:
		return "!"

	// Bitwise
	case AMPERSAND:
		return "&"
	case PIPE:
		return "|"
	case TILDE:
		return "~"
	case CARET:
		return "^"
	case LESS_LESS:
		return "<<"
	case GREATER_GREATER:
		return ">>"

	// Math
	case PLUS:
		return "+"
	case MINUS:
		return "-"
	case STAR:
		return "*"
	case STAR_STAR:
		return "**"
	case SLASH:
		return "/"
	case PERCENT:
		return "%"

	// Assignment
	case EQUAL:
		return "="
	case PLUS_EQUAL:
		return "+="
	case MINUS_EQUAL:
		return "-="
	case STAR_EQUAL:
		return "*="
	case STAR_STAR_EQUAL:
		return "**="
	case SLASH_EQUAL:
		return "/="
	case PERCENT_EQUAL:
		return "%="
	case LESS_LESS_EQUAL:
		return "<<="
	case GREATER_GREATER_EQUAL:
		return ">>="
	case AMPERSAND_EQUAL:
		return "&="
	case PIPE_EQUAL:
		return "|="
	case CARET_EQUAL:
		return "^="

	// Control flow
	case IF:
		return "if"
	case ELIF:
		return "elif"
	case ELSE:
		return "else"
	case FOR:
		return "for"
	case WHILE:
		return "while"
	case BREAK:
		return "break"
	case CONTINUE:
		return "continue"
	case PASS:
		return "pass"
	case RETURN:
		return "return"
	case MATCH:
		return "match"
	case WHEN:
		return "when"

	// Keywords
	case AS:
		return "as"
	case CLASS:
		return "class"
	case CONST:
		return "const"
	case ENUM:
		return "enum"
	case EXTENDS:
		return "extends"
	case FUNCTION:
		return "fn"
	case IMPORT:
		return "import"
	case IN:
		return "in"
	case IS:
		return "is"
	case MOD:
		return "mod"
	case SELF:
		return "self"
	case SIGNAL:
		return "signal"
	case TRAIT:
		return "trait"
	case TYPE:
		return "type"
	case USES:
		return "uses"
	case VAR:
		return "var"
	case VOID:
		return "void"

	// Punctuation
	case BRACKET_OPEN:
		return "["
	case BRACKET_CLOSE:
		return "]"
	case BRACE_OPEN:
		return "{"
	case BRACE_CLOSE:
		return "}"
	case PARENTHESIS_OPEN:
		return "("
	case PARENTHESIS_CLOSE:
		return ")"
	case COMMA:
		return ","
	case SEMICOLON:
		return ";"
	case PERIOD:
		return "."
	case PERIOD_PERIOD:
		return ".."
	case PERIOD_PERIOD_PERIOD:
		return "..."
	case COLON:
		return ":"
	case DOLLAR:
		return "$"
	case FORWARD_ARROW:
		return "->"
	case UNDERSCORE:
		return "_"

	// Whitespace
	case NEWLINE:
		return "Newline"

	// Constants
	case CONST_PI:
		return "PI"
	case CONST_TAU:
		return "TAU"
	case CONST_INF:
		return "INF"
	case CONST_NAN:
		return "NaN"

	// Error message improvement
	case VCS_CONFLICT_MARKER:
		return "VCS conflict marker"
	case BACKTICK:
		return "`"
	case QUESTION_MARK:
		return "?"

	// Special
	case ERROR:
		return "Error"
	case EOF:
		return "End of file"

	default:
		panic("Using token type out of the enum.")
	}
}


func (t *Token) GetDebugName() string {
	if t.Type == IDENTIFIER {
		return "identifier: " + string(t.Source)
	}

	if t.Type == LITERAL {
		return "Literal: " + string(t.Source)
	}

	if t.Type == ERROR {
		s, ok := t.Literal.(string)
		if !ok {
			return "Error"
		} else {
			return "Error: " + s
		}
	}
	return t.GetName()
}

func (t *Token) CanPrecedeBinOP() bool {
	switch t.Type {
	case IDENTIFIER, LITERAL, SELF, BRACKET_CLOSE,
		BRACE_CLOSE, PARENTHESIS_CLOSE,
		CONST_PI, CONST_TAU, CONST_INF, CONST_NAN:
		return true
	default:
		return false
	}
}

func (t *Token) IsIdentifier() bool {
	switch t.Type {
	case IDENTIFIER, MATCH, WHEN, CONST_PI,
		CONST_TAU, CONST_INF, CONST_NAN:
		return true
	default:
		return false
	}
}

func (t *Token) IsNodeName() bool {
	switch t.Type {
	case IDENTIFIER, AND, AS, BREAK,
		CLASS, CONST, CONST_PI, CONST_INF,
		CONST_NAN, CONST_TAU, CONTINUE, ELIF, ELSE,
		ENUM, EXTENDS, FOR, FUNCTION, IF, IMPORT, IN, IS, MATCH, MOD,
		NOT, OR, PASS, RETURN, SELF, SIGNAL, TRAIT, TYPE,
		UNDERSCORE, USES, VAR, VOID, WHILE, WHEN:
		return true
	default:
		return false
	}
}
