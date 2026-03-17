package output

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"
)

// Format represents the output format.
type Format string

const (
	FormatTable Format = "table"
	FormatJSON  Format = "json"
)

// Printer handles formatted output.
type Printer struct {
	format  Format
	noColor bool
	out     io.Writer
}

// NewPrinter creates a new output printer.
func NewPrinter(format Format, noColor bool) *Printer {
	return &Printer{
		format:  format,
		noColor: noColor,
		out:     os.Stdout,
	}
}

// PrintJSON outputs data as formatted JSON.
func (p *Printer) PrintJSON(v interface{}) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	fmt.Fprintln(p.out, string(data))
	return nil
}

// Table creates a new table builder.
func (p *Printer) Table(headers ...string) *TableBuilder {
	return &TableBuilder{
		printer: p,
		headers: headers,
	}
}

// PrintResource prints a single resource in the chosen format.
// For table format, it prints key-value pairs.
func (p *Printer) PrintResource(fields []KeyValue) error {
	if p.format == FormatJSON {
		m := make(map[string]interface{})
		for _, f := range fields {
			m[f.Key] = f.Value
		}
		return p.PrintJSON(m)
	}

	w := tabwriter.NewWriter(p.out, 0, 4, 2, ' ', 0)
	for _, f := range fields {
		fmt.Fprintf(w, "%s:\t%v\n", f.Key, f.Value)
	}
	return w.Flush()
}

// Success prints a success message.
func (p *Printer) Success(msg string, args ...interface{}) {
	prefix := "✓"
	if p.noColor {
		prefix = "OK"
	}
	fmt.Fprintf(p.out, "%s %s\n", prefix, fmt.Sprintf(msg, args...))
}

// Error prints an error message to stderr.
func (p *Printer) Error(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "Error: %s\n", fmt.Sprintf(msg, args...))
}

// KeyValue represents a key-value pair for resource display.
type KeyValue struct {
	Key   string
	Value interface{}
}

// TableBuilder builds a table output.
type TableBuilder struct {
	printer *Printer
	headers []string
	rows    [][]string
}

// AddRow adds a row to the table.
func (t *TableBuilder) AddRow(values ...string) *TableBuilder {
	t.rows = append(t.rows, values)
	return t
}

// Print outputs the table.
func (t *TableBuilder) Print() error {
	if t.printer.format == FormatJSON {
		items := make([]map[string]string, len(t.rows))
		for i, row := range t.rows {
			item := make(map[string]string)
			for j, h := range t.headers {
				if j < len(row) {
					item[strings.ToLower(h)] = row[j]
				}
			}
			items[i] = item
		}
		return t.printer.PrintJSON(items)
	}

	w := tabwriter.NewWriter(t.printer.out, 0, 4, 2, ' ', 0)
	fmt.Fprintln(w, strings.Join(t.headers, "\t"))
	fmt.Fprintln(w, strings.Repeat("-\t", len(t.headers)))
	for _, row := range t.rows {
		fmt.Fprintln(w, strings.Join(row, "\t"))
	}
	return w.Flush()
}
