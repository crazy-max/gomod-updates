package gomodupdates

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

func tableHeaders() []string {
	return []string{"Module", "Version", "New Version", "Direct", "Valid Timestamps"}
}

func tableValues(row Row) []string {
	return []string{
		row.Module,
		row.Version,
		row.NewVersion,
		strconv.FormatBool(row.Direct),
		strconv.FormatBool(row.ValidTimestamps),
	}
}

// RenderTable writes rows in the selected table format.
func RenderTable(w io.Writer, rows []Row, opts Options) {
	headers := tableHeaders()
	values := make([][]string, 0, len(rows))
	widths := make([]int, len(headers))
	for i, header := range headers {
		widths[i] = len(header)
	}
	for _, row := range rows {
		value := tableValues(row)
		values = append(values, value)
		for i := range value {
			if len(value[i]) > widths[i] {
				widths[i] = len(value[i])
			}
		}
	}

	if opts.Format == FormatMarkdown {
		renderMarkdown(w, headers, values)
		return
	}
	renderDefault(w, headers, values, widths)
}

func renderDefault(w io.Writer, headers []string, values [][]string, widths []int) {
	border := func() {
		_, _ = fmt.Fprint(w, "+")
		for _, width := range widths {
			_, _ = fmt.Fprint(w, strings.Repeat("-", width+2), "+")
		}
		_, _ = fmt.Fprintln(w)
	}
	row := func(values []string) {
		_, _ = fmt.Fprint(w, "|")
		for i, value := range values {
			_, _ = fmt.Fprintf(w, " %-*s |", widths[i], value)
		}
		_, _ = fmt.Fprintln(w)
	}

	border()
	row(headers)
	border()
	for _, value := range values {
		row(value)
	}
	border()
}

func renderMarkdown(w io.Writer, headers []string, values [][]string) {
	writeMarkdownRow(w, headers)
	separators := make([]string, len(headers))
	for i := range separators {
		separators[i] = "---"
	}
	writeMarkdownRow(w, separators)
	for _, value := range values {
		writeMarkdownRow(w, value)
	}
}

func writeMarkdownRow(w io.Writer, values []string) {
	_, _ = fmt.Fprint(w, "|")
	for _, value := range values {
		_, _ = fmt.Fprintf(w, " %s |", value)
	}
	_, _ = fmt.Fprintln(w)
}
