package gomodupdates

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"
)

const moduleStream = `{
	"Path": "github.com/example/root",
	"Main": true
}
{
	"Path": "github.com/example/direct",
	"Version": "v1.0.0",
	"Update": {
		"Path": "github.com/example/direct",
		"Version": "v1.1.0"
	}
}
{
	"Path": "github.com/example/major",
	"Version": "v1.0.0"
}
{
	"Path": "github.com/example/indirect",
	"Version": "v1.0.0",
	"Indirect": true,
	"Update": {
		"Path": "github.com/example/indirect",
		"Version": "v1.1.0"
	}
}`

func TestRunRendersMajorColumn(t *testing.T) {
	var out bytes.Buffer
	err := Run(context.Background(), strings.NewReader(moduleStream), &out, Options{
		Update: true,
		Direct: true,
		Major:  true,
		Lister: fakeVersionLister{
			"github.com/example/direct/v2": {"v2.0.0"},
			"github.com/example/major/v2":  {"v2.3.4"},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	got := out.String()
	for _, want := range []string{
		"github.com/example/direct/v2@v2.0.0",
		"github.com/example/major/v2@v2.3.4",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("expected output to contain %q:\n%s", want, got)
		}
	}
	if strings.Contains(got, "New Major Version") {
		t.Fatalf("expected major updates to use New Version column:\n%s", got)
	}
	if strings.Contains(got, "github.com/example/indirect") {
		t.Fatalf("expected indirect module to be filtered:\n%s", got)
	}
}

func TestRunCIReturnsErrOutdated(t *testing.T) {
	err := Run(context.Background(), strings.NewReader(moduleStream), ioDiscard{}, Options{
		Update: true,
		CI:     true,
	})
	if !errors.Is(err, ErrOutdated) {
		t.Fatalf("expected ErrOutdated, got %v", err)
	}
}

type ioDiscard struct{}

func (ioDiscard) Write(p []byte) (int, error) {
	return len(p), nil
}
