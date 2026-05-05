package gomodupdates

import (
	"bytes"
	"context"
	"os/exec"
)

// GoListModules runs go list -u -m -json all and returns the JSON stream.
func GoListModules(ctx context.Context, goCommand, mod string) ([]byte, error) {
	if goCommand == "" {
		goCommand = "go"
	}

	args := []string{"list"}
	if mod != "" {
		args = append(args, "-mod="+mod)
	}
	args = append(args, "-u", "-m", "-json", "all")

	cmd := exec.CommandContext(ctx, goCommand, args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil && stderr.Len() > 0 {
		return nil, &GoListError{Err: err, Stderr: stderr.String()}
	}
	return out, err
}

// GoListError wraps go list failures that include stderr.
type GoListError struct {
	Err    error
	Stderr string
}

func (e *GoListError) Error() string {
	return e.Err.Error() + ": " + e.Stderr
}

func (e *GoListError) Unwrap() error {
	return e.Err
}
