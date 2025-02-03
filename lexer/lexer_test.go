package lexer_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/tobiashort/gox/lexer"
	"github.com/tobiashort/gox/lexer/assert"
)

func TestTokenizeExamples(t *testing.T) {
	entries, err := os.ReadDir(filepath.Join("..", "examples"))
	assert.Nil(err)
	for _, entry := range entries {
		t.Run(entry.Name(), func(*testing.T) {
			data, err := os.ReadFile(filepath.Join("..", "examples", entry.Name()))
			assert.Nil(err)
			source := string(data)
			lexer := lexer.NewLexer()
			lexer.Tokenize(source)
		})
	}
}
