package gomodupdates

import (
	"bytes"
	"strings"
	"testing"
)

func TestRenderMarkdown(t *testing.T) {
	var out bytes.Buffer
	RenderTable(&out, []Row{
		{
			Module:          "github.com/example/mod",
			Version:         "v1.0.0",
			NewVersion:      "github.com/example/mod/v2@v2.0.0",
			Direct:          true,
			ValidTimestamps: true,
		},
	}, Options{Major: true, Format: FormatMarkdown})

	got := out.String()
	if !strings.Contains(got, "| Module | Version | New Version | Direct | Valid Timestamps |") {
		t.Fatalf("unexpected markdown table:\n%s", got)
	}
	if !strings.Contains(got, "| github.com/example/mod | v1.0.0 | github.com/example/mod/v2@v2.0.0 | true | true |") {
		t.Fatalf("unexpected markdown table:\n%s", got)
	}
}
