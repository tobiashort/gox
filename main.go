package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/sanity-io/litter"

	"github.com/tobiashort/gox/assert"
	"github.com/tobiashort/gox/ast"
	"github.com/tobiashort/gox/lexer"
	"github.com/tobiashort/gox/parser"
	"github.com/tobiashort/gox/transpiler"
)

func usage() {
	fmt.Fprintf(os.Stderr, `Usage:
	gox tokenize FILE
	gox parse FILE
	gox transpile FILE
	gox run FILE
`)
}

func tokenize(file string) []lexer.Token {
	data, err := os.ReadFile(file)
	assert.Nil(err)
	lexer := lexer.NewLexer()
	lexer.Tokenize(string(data))
	return lexer.Tokens
}

func parse(file string) []ast.Stmt {
	data, err := os.ReadFile(file)
	assert.Nil(err)
	parser := parser.NewParser()
	parser.Parse(string(data))
	return parser.Stmts
}

func transpile(file string) string {
	data, err := os.ReadFile(file)
	assert.Nil(err)
	transpiler := transpiler.NewTranspiler()
	transpiler.Transpile(string(data))
	return transpiler.String()
}

func main() {
	flag.Usage = usage
	flag.Parse()

	task := ""
	if flag.NArg() > 0 {
		task = flag.Arg(0)
	}

	switch task {
	case "tokenize":
		if flag.NArg() != 2 {
			fmt.Fprintln(os.Stderr, "must provide file")
			os.Exit(1)
		}
		file := flag.Arg(1)
		tokens := tokenize(file)
		for _, token := range tokens {
			fmt.Println(token)
		}
	case "parse":
		if flag.NArg() != 2 {
			fmt.Fprintln(os.Stderr, "must provide file")
			os.Exit(1)
		}
		file := flag.Arg(1)
		ast := parse(file)
		litter.Dump(ast)
	case "transpile":
		if flag.NArg() != 2 {
			fmt.Fprintln(os.Stderr, "must provide file")
			os.Exit(1)
		}
		file := flag.Arg(1)
		source := transpile(file)
		fmt.Println(source)
	case "run":
		if flag.NArg() != 2 {
			fmt.Fprintln(os.Stderr, "must provide file")
			os.Exit(1)
		}
		file := flag.Arg(1)
		source := transpile(file)
		tempDir, err := os.MkdirTemp(os.TempDir(), "gox")
		assert.Nil(err)
		defer os.RemoveAll(tempDir)
		goFile, err := os.OpenFile(filepath.Join(tempDir, strings.TrimSuffix(filepath.Base(file), ".gox")+".go"), os.O_CREATE|os.O_RDWR, 0o644)
		fmt.Println(goFile.Name())
		assert.Nil(err)
		_, err = io.Copy(goFile, bytes.NewReader([]byte(source)))
		assert.Nil(err)
		cmd := exec.Command("go", "run", goFile.Name())
		output, err := cmd.CombinedOutput()
		fmt.Print(string(output))
		if err != nil {
			os.Exit(1)
		}
	default:
		usage()
		os.Exit(1)
	}
}
