package gomodupdates

import (
	"context"
	"encoding/json"
	"errors"
	"io"

	"golang.org/x/sync/errgroup"
)

// ErrOutdated is returned when CI mode is enabled and at least one update is found.
var ErrOutdated = errors.New("outdated modules found")

// OutputFormat specifies the supported table rendering formats.
type OutputFormat string

const (
	// FormatDefault is the default bordered table output.
	FormatDefault OutputFormat = "default"
	// FormatMarkdown is a Markdown table.
	FormatMarkdown OutputFormat = "markdown"
)

const defaultMajorConcurrency = 8

// Options configures module filtering and rendering.
type Options struct {
	Update           bool
	Direct           bool
	Major            bool
	CI               bool
	Format           OutputFormat
	MajorLimit       int
	MajorConcurrency int
	Lister           VersionLister
}

// Run converts go list -u -m -json all output into a table.
func Run(ctx context.Context, in io.Reader, out io.Writer, opts Options) error {
	modules, err := DecodeModules(in)
	if err != nil {
		return err
	}

	rows, err := Rows(ctx, modules, opts)
	if err != nil {
		return err
	}
	if len(rows) > 0 {
		RenderTable(out, rows, opts)
	}
	if opts.CI && hasOutdated(rows) {
		return ErrOutdated
	}
	return nil
}

// DecodeModules decodes the JSON stream returned by go list -m -json.
func DecodeModules(in io.Reader) ([]Module, error) {
	var modules []Module
	dec := json.NewDecoder(in)
	for {
		var mod Module
		err := dec.Decode(&mod)
		if err != nil {
			if errors.Is(err, io.EOF) {
				return modules, nil
			}
			return nil, err
		}
		modules = append(modules, mod)
	}
}

// Rows filters modules and resolves optional major update candidates.
func Rows(ctx context.Context, modules []Module, opts Options) ([]Row, error) {
	rows := make([]Row, 0, len(modules))
	majorIndexes := make([]int, 0, len(modules))
	for _, mod := range modules {
		if mod.Main {
			continue
		}
		if opts.Direct && mod.Indirect {
			continue
		}

		rows = append(rows, Row{
			Module:          mod.Path,
			Version:         mod.CurrentVersion(),
			NewVersion:      mod.NewVersion(),
			Direct:          !mod.Indirect,
			ValidTimestamps: !mod.InvalidTimestamp(),
		})
		if opts.Major {
			majorIndexes = append(majorIndexes, len(rows)-1)
		}
	}

	if len(majorIndexes) > 0 {
		if err := resolveMajorUpdates(ctx, rows, majorIndexes, opts); err != nil {
			return nil, err
		}
	}

	if opts.Update {
		filtered := rows[:0]
		for _, row := range rows {
			if row.HasUpdate() {
				filtered = append(filtered, row)
			}
		}
		rows = filtered
	}

	return rows, nil
}

func resolveMajorUpdates(ctx context.Context, rows []Row, indexes []int, opts Options) error {
	eg, ctx := errgroup.WithContext(ctx)
	concurrency := opts.MajorConcurrency
	if concurrency == 0 {
		concurrency = defaultMajorConcurrency
	}
	if concurrency > 0 {
		eg.SetLimit(concurrency)
	}

	for _, index := range indexes {
		eg.Go(func() error {
			major, err := FindMajorUpdate(ctx, rows[index].Module, rows[index].Version, opts.Lister, opts.MajorLimit)
			if err != nil {
				return err
			}
			if version := major.String(); version != "" {
				rows[index].NewVersion = version
			}
			return nil
		})
	}

	return eg.Wait()
}

func hasOutdated(rows []Row) bool {
	for _, row := range rows {
		if row.HasUpdate() {
			return true
		}
	}
	return false
}
