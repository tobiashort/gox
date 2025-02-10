package lexer

import "fmt"

type TokenType = string

const (
	TokenString     = "STRING"
	TokenNumber     = "NUMBER"
	TokenIdentifier = "IDENTIFIER"

	// operators
	TokenAssign     = "ASSIGN"
	TokenDeclAssign = "DECL_ASSIGN"
	TokenPlus       = "PLUS"
	TokenStar       = "STAR"

	//  punctuation
	TokenDot        = "DOT"
	TokenParenOpen  = "PAREN_OPEN"
	TokenParenClose = "PAREN_CLOSE"
	TokenBraceOpen  = "BRACE_OPEN"
	TokenBraceClose = "BRACE_CLOSE"
	TokenNewLine    = "NEW_LINE"
	TokenComma      = "COMMA"

	// Keywords
	TokenPackage = "PACKAGE"
	TokenImport  = "IMPORT"
	TokenFunc    = "FUNC"
	TokenVar     = "VAR"
	TokenOrPanic = "OR_PANIC"
	TokenReturn  = "RETURN"

	TokenEOF = "EOF"
)

var Keywords = map[string]TokenType{
	"package":  TokenPackage,
	"import":   TokenImport,
	"func":     TokenFunc,
	"var":      TokenVar,
	"or_panic": TokenOrPanic,
	"return":   TokenReturn,
}

func IsKeyword(value string) bool {
	_, exists := Keywords[value]
	return exists
}

type Token struct {
	Type   string
	Value  string
	Line   int
	Column int
}

func NewToken(_type string, value string, line, column int) Token {
	return Token{
		Type:   _type,
		Value:  value,
		Line:   line,
		Column: column,
	}
}

func (t Token) String() string {
	if t.Value == "" {
		return fmt.Sprintf("%s", t.Type)
	}
	return fmt.Sprintf("%s (%s)", t.Type, t.Value)
}
