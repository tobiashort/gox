package lexer

import "regexp"

type PatternHandler = func(lexer *Lexer, regex *regexp.Regexp)

type Pattern struct {
	Regex   *regexp.Regexp
	Handler PatternHandler
}

func SkipHandler() PatternHandler {
	return func(lexer *Lexer, regex *regexp.Regexp) {
		loc := regex.FindStringIndex(lexer.Remainder())
		lexer.Pos += loc[1]
	}
}

func DefaultHandler(_type TokenType) PatternHandler {
	return func(lexer *Lexer, regex *regexp.Regexp) {
		loc := regex.FindStringIndex(lexer.Remainder())
		lexer.Add(NewToken(_type, "", lexer.Line(), lexer.Column()))
		lexer.Pos += loc[1]
	}
}

func IdentifierHandler() PatternHandler {
	return func(lexer *Lexer, regex *regexp.Regexp) {
		value := regex.FindString(lexer.Remainder())
		if IsKeyword(value) {
			lexer.Add(NewToken(Keywords[value], "", lexer.Line(), lexer.Column()))
		} else {
			lexer.Add(NewToken(TokenIdentifier, value, lexer.Line(), lexer.Column()))
		}
		lexer.Pos += len(value)
	}
}

func StringHandler() PatternHandler {
	return func(lexer *Lexer, regex *regexp.Regexp) {
		loc := regex.FindStringIndex(lexer.Remainder())
		lexer.Add(NewToken(TokenString, lexer.Remainder()[loc[0]+1:loc[1]-1], lexer.Line(), lexer.Column()))
		lexer.Pos += loc[1]
	}
}

func NumberHandler() PatternHandler {
	return func(lexer *Lexer, regex *regexp.Regexp) {
		value := regex.FindString(lexer.Remainder())
		lexer.Add(NewToken(TokenNumber, value, lexer.Line(), lexer.Column()))
		lexer.Pos += len(value)
	}
}

var Patterns = []Pattern{
	{regexp.MustCompile("^\\n+"), DefaultHandler(TokenNewLine)},
	{regexp.MustCompile("^\\s+"), SkipHandler()},
	{regexp.MustCompile("^[A-Za-z_][A-Za-z0-9_]*"), IdentifierHandler()},
	{regexp.MustCompile(`^"[^"]*"`), StringHandler()},
	{regexp.MustCompile("^\\d+(\\.\\d+)?"), NumberHandler()},
	{regexp.MustCompile("^package"), DefaultHandler(TokenPackage)},
	{regexp.MustCompile("^import"), DefaultHandler(TokenImport)},
	{regexp.MustCompile("^="), DefaultHandler(TokenAssign)},
	{regexp.MustCompile("^:="), DefaultHandler(TokenDeclAssign)},
	{regexp.MustCompile("^\\+"), DefaultHandler(TokenPlus)},
	{regexp.MustCompile("^\\*"), DefaultHandler(TokenStar)},
	{regexp.MustCompile("^\\("), DefaultHandler(TokenParenOpen)},
	{regexp.MustCompile("^\\)"), DefaultHandler(TokenParenClose)},
	{regexp.MustCompile("^\\{"), DefaultHandler(TokenBraceOpen)},
	{regexp.MustCompile("^\\}"), DefaultHandler(TokenBraceClose)},
	{regexp.MustCompile("^\\."), DefaultHandler(TokenDot)},
	{regexp.MustCompile("^,"), DefaultHandler(TokenComma)},
}
