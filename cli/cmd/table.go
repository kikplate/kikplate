package cmd

import (
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
)

type Table struct {
	headers []string
	rows    [][]string
}

func NewTable(headers ...string) *Table {
	return &Table{headers: headers}
}

func (t *Table) Row(values ...string) *Table {
	t.rows = append(t.rows, values)
	return t
}

func (t *Table) Print() {
	tw := table.NewWriter()
	tw.SetOutputMirror(os.Stdout)

	headerRow := make(table.Row, len(t.headers))
	for i, h := range t.headers {
		headerRow[i] = h
	}
	tw.AppendHeader(headerRow)

	for _, row := range t.rows {
		r := make(table.Row, len(row))
		for i, v := range row {
			r[i] = v
		}
		tw.AppendRow(r)
	}

	tw.SetStyle(table.StyleLight)
	tw.Render()
}

func (t *Table) PrintWithTotal(total int) {
	t.Print()
	fmt.Printf("\nTotal: %d\n", total)
}

func truncate(s string, max int) string {
	if len(s) > max {
		return s[:max-3] + "..."
	}
	return s
}

func boolYesNo(b bool) string {
	if b {
		return "yes"
	}
	return "no"
}

func deref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
