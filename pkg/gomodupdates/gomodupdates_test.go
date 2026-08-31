package gomodupdates

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestToolPathsFromGoMod(t *testing.T) {
	path := filepath.Join(t.TempDir(), "go.mod")
	err := os.WriteFile(path, []byte(`module github.com/example/root

go 1.26.0

require github.com/example/tool v1.0.0 // indirect

tool github.com/example/tool/cmd/tool
`), 0644)
	if err != nil {
		t.Fatal(err)
	}

	got, err := ToolPathsFromGoMod(path)
	if err != nil {
		t.Fatal(err)
	}

	want := []string{"github.com/example/tool/cmd/tool"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected tool paths %v, got %v", want, got)
	}
}
