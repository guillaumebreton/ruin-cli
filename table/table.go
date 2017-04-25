package table

import (
	"fmt"
	"io"
	"strings"
)

type Data struct {
	data []string
}
type Table struct {
	Rows      []Row
	Header    Data
	Separator string
	Border    bool
	Align     bool
}

type Row interface {
	render(io.Writer, []int)
}

func (d Data) render(writer io.Writer, constraints []int) {
	arr := make([]string, len(constraints))
	for k, c := range d.data {
		width := constraints[k]
		format := fmt.Sprintf(" %% %ds ", width)
		arr[k] = fmt.Sprintf(format, c)
	}
	render(writer, arr, "|")
}

type Separator struct {
}

func (s Separator) render(writer io.Writer, constraints []int) {
	arr := make([]string, len(constraints))
	for k, c := range constraints {
		s := "-"
		for i := 0; i < c; i++ {
			s += "-"
		}
		s += "-"
		arr[k] = s
	}
	render(writer, arr, "+")
}

func render(writer io.Writer, data []string, separator string) {
	s := separator + strings.Join(data, separator) + separator + "\n"
	writer.Write([]byte(s))
}

func NewTable() *Table {
	return &Table{
		Rows:   make([]Row, 0),
		Header: Data{make([]string, 0)},
	}
}

func (t *Table) SetHeader(headers []string) {
	t.Header = Data{headers}
}
func (t *Table) Append(data []string) {
	t.Rows = append(t.Rows, Data{data})
}
func (t *Table) AppendSeparator() {
	t.Rows = append(t.Rows, Separator{})
}

func (t *Table) Render(writer io.Writer) {
	s := Separator{}
	//Get the maximum length of each cell
	widths := t.computeCellWith()
	s.render(writer, widths)
	t.Header.render(writer, widths)
	s.render(writer, widths)
	t.AppendSeparator()
	previousIsSeparator := true
	for _, r := range t.Rows {
		_, isSeparator := r.(Separator)
		if !isSeparator {
			r.render(writer, widths)
		} else if !previousIsSeparator {
			r.render(writer, widths)
		}
		previousIsSeparator = isSeparator
	}
	t.AppendSeparator()

}

func (t *Table) computeCellWith() []int {
	max := 0
	for _, v := range t.Rows {
		if r, ok := v.(Data); ok {
			if len(r.data) > max {
				max = len(r.data)
			}
		}
	}
	widths := make([]int, max)
	for _, v := range t.Rows {
		if r, ok := v.(Data); ok {
			for j, s := range r.data {
				if len(s) > widths[j] {
					widths[j] = len(s)
				}
			}
		}
	}
	return widths
}
