package main

import (
	"context"
	"errors"
	"io"
	"os"
	"testing"
)

func TestInputUsesNonEmptyStdin(t *testing.T) {
	const payload = `{"Path":"github.com/example/root","Main":true}`
	stdin := stdinFile(t, payload)
	called := false

	in, err := input(context.Background(), stdin, "mod", func(context.Context, string, string) ([]byte, error) {
		called = true
		return nil, errors.New("unexpected go list call")
	})
	if err != nil {
		t.Fatal(err)
	}
	if called {
		t.Fatal("expected non-empty stdin to be used directly")
	}

	got, err := io.ReadAll(in)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != payload {
		t.Fatalf("expected stdin payload %q, got %q", payload, string(got))
	}
}

func TestInputFallsBackWhenNonTTYStdinIsEmpty(t *testing.T) {
	const payload = `{"Path":"github.com/example/root","Main":true}`
	stdin := stdinFile(t, "")

	in, err := input(context.Background(), stdin, "readonly", func(_ context.Context, goCommand, mod string) ([]byte, error) {
		if goCommand != "" {
			t.Fatalf("expected default go command, got %q", goCommand)
		}
		if mod != "readonly" {
			t.Fatalf("expected mod mode readonly, got %q", mod)
		}
		return []byte(payload), nil
	})
	if err != nil {
		t.Fatal(err)
	}

	got, err := io.ReadAll(in)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != payload {
		t.Fatalf("expected go list payload %q, got %q", payload, string(got))
	}
}

func stdinFile(t *testing.T, contents string) *os.File {
	t.Helper()

	f, err := os.CreateTemp(t.TempDir(), "stdin")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(contents); err != nil {
		t.Fatal(err)
	}
	if _, err := f.Seek(0, io.SeekStart); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := f.Close(); err != nil {
			t.Fatal(err)
		}
	})
	return f
}
