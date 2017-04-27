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
	Rows            []Row
	Header          Data
	CenterSeparator string
	ColumnSeparator string
	RowSeparator    string
	Border          bool
}

type Row interface {
	render(io.Writer, []int, *Table)
}

func (d Data) render(writer io.Writer, constraints []int, t *Table) {
	arr := make([]string, len(constraints))
	for k, c := range d.data {
		width := constraints[k]
		var format string
		if t.Border && k == 0 {
			format = fmt.Sprintf("%% %ds ", width)
		} else if k == len(d.data)-1 {

			format = fmt.Sprintf(" %% %ds", width)
		} else {
			format = fmt.Sprintf(" %% %ds ", width)

		}
		arr[k] = fmt.Sprintf(format, c)
	}
	t.render(writer, arr, t.ColumnSeparator)
}

type Separator struct {
}

func (s Separator) render(writer io.Writer, constraints []int, t *Table) {
	arr := make([]string, len(constraints))
	for k, c := range constraints {
		s := t.RowSeparator
		for i := 0; i < c; i++ {
			s += t.RowSeparator
		}
		s += t.RowSeparator
		arr[k] = s
	}
	t.render(writer, arr, t.CenterSeparator)
}

func (t *Table) render(writer io.Writer, data []string, separator string) {
	var s string
	if t.Border {
		s = separator + strings.Join(data, separator) + separator + "\n"
	} else {
		s = strings.Join(data, separator) + "\n"
	}
	writer.Write([]byte(s))
}

func NewTable() *Table {
	return &Table{
		Rows:            make([]Row, 0),
		Header:          Data{make([]string, 0)},
		ColumnSeparator: "|",
		RowSeparator:    "-",
		CenterSeparator: "+",
		Border:          true,
	}
}

func (t *Table) SetBorder(b bool) {
	t.Border = b
}

func (t *Table) SetCenterSeparator(s string) {
	t.CenterSeparator = s
}

func (t *Table) SetColumnSeparator(s string) {
	t.ColumnSeparator = s
}

func (t *Table) SetRowSeparator(s string) {
	t.RowSeparator = s
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
	if t.Border {
		s.render(writer, widths, t)
	}
	if len(t.Header.data) > 0 {
		t.Header.render(writer, widths, t)
		s.render(writer, widths, t)
	}
	previousIsSeparator := true
	for _, r := range t.Rows {
		_, isSeparator := r.(Separator)
		if !isSeparator {
			r.render(writer, widths, t)
		} else if !previousIsSeparator {
			r.render(writer, widths, t)
		}
		previousIsSeparator = isSeparator
	}
	if t.Border {
		s.render(writer, widths, t)
	}
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
