package transpiler_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/tobiashort/gox/assert"
	"github.com/tobiashort/gox/transpiler"
)

func TestTranspileExamples(t *testing.T) {
	entries, err := os.ReadDir(filepath.Join("..", "examples"))
	assert.Nil(err)
	for _, entry := range entries {
		t.Run(entry.Name(), func(*testing.T) {
			data, err := os.ReadFile(filepath.Join("..", "examples", entry.Name()))
			assert.Nil(err)
			source := string(data)
			transpiler := transpiler.NewTranspiler()
			transpiler.Transpile(source)
			t.Log(transpiler.String())
		})
	}
}
