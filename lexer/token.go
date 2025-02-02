package lexer

import "fmt"

type TokenType = string

const (
	TokenString     = "STRING"
	TokenIdentifier = "IDENTIFIER"

	//  punctuation
	TokenDot        = "DOT"
	TokenParenOpen  = "PAREN_OPEN"
	TokenParenClose = "PAREN_CLOSE"
	TokenBraceOpen  = "BRACE_CLOSE"
	TokenBraceClose = "BRACE_CLOSE"

	// Keywords
	TokenPackage = "PACKAGE"
	TokenImport  = "IMPORT"
	TokenFunc    = "FUNC"
	TokenOrPanic = "OR_PANIC"
	TokenReturn  = "RETURN"

	TokenEOF = "EOF"
)

var Keywords = map[string]TokenType{
	"package":  TokenPackage,
	"import":   TokenImport,
	"func":     TokenFunc,
	"or_panic": TokenOrPanic,
	"return":   TokenReturn,
}

func IsKeyword(value string) bool {
	_, exists := Keywords[value]
	return exists
}

type Token struct {
	Type  string
	Value any
}

func NewToken(_type string, value any) Token {
	return Token{
		Type:  _type,
		Value: value,
	}
}

func (t Token) String() string {
	if t.Value == nil {
		return fmt.Sprintf("%s", t.Type)
	}
	return fmt.Sprintf("%s (%s)", t.Type, t.Value)
}
