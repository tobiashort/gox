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
		lexer.Add(NewToken(_type, nil))
		lexer.Pos += loc[1]
	}
}

func IdentifierHandler() PatternHandler {
	return func(lexer *Lexer, regex *regexp.Regexp) {
		value := regex.FindString(lexer.Remainder())
		if IsKeyword(value) {
			lexer.Add(NewToken(Keywords[value], nil))
		} else {
			lexer.Add(NewToken(TokenIdentifier, value))
		}
		lexer.Pos += len(value)
	}
}

func StringHandler() PatternHandler {
	return func(lexer *Lexer, regex *regexp.Regexp) {
		loc := regex.FindStringIndex(lexer.Remainder())
		lexer.Add(NewToken(TokenString, lexer.Remainder()[loc[0]+1:loc[1]-1]))
		lexer.Pos += loc[1]
	}
}

var Patterns = []Pattern{
	{regexp.MustCompile("^\\s+"), SkipHandler()},
	{regexp.MustCompile("^[A-Za-z_][A-Za-z0-9_]*"), IdentifierHandler()},
	{regexp.MustCompile(`^"[^"]*"`), StringHandler()},
	{regexp.MustCompile("^package"), DefaultHandler(TokenPackage)},
	{regexp.MustCompile("^import"), DefaultHandler(TokenImport)},
	{regexp.MustCompile("^\\("), DefaultHandler(TokenParenOpen)},
	{regexp.MustCompile("^\\)"), DefaultHandler(TokenParenClose)},
	{regexp.MustCompile("^\\{"), DefaultHandler(TokenBraceOpen)},
	{regexp.MustCompile("^\\}"), DefaultHandler(TokenBraceClose)},
	{regexp.MustCompile("^\\."), DefaultHandler(TokenDot)},
}
