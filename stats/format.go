package stats

import (
	"fmt"

	"github.com/alexeyco/simpletable"
)

// TableDatum generic type for table datum
type TableDatum []interface{}

// TableData generic type for table data
type TableData []TableDatum

// Table generic "class" for table data
type Table struct {
	Data  TableData
	Table *simpletable.Table
}

func buildRow(index int, data interface{}) *simpletable.Cell {
	cell := &simpletable.Cell{
		Align: simpletable.AlignLeft,
	}
	switch data.(type) {
	case string:
		cell.Text = data.(string)
	case int:
		cell.Text = fmt.Sprintf("%d", data.(int))
		cell.Align = simpletable.AlignRight
	case float64:
		cell.Text = fmt.Sprintf("%.2f%%", data.(float64))
		cell.Align = simpletable.AlignLeft
	}
	return cell
}

func (t *Table) SetCellsForData(transform func(int, interface{}) *simpletable.Cell) {
	rows := make([][]*simpletable.Cell, 0)
	for _, val := range t.Data {
		row := make([]*simpletable.Cell, 0)
		for i, data := range val {
			if transform == nil {
				row = append(row, buildRow(i, data))
			} else {
				row = append(row, transform(i, data))
			}
		}
		rows = append(rows, row)
	}
	t.Table.Body.Cells = rows
}

func (t *Table) AddHeader(keys []string) {
	cells := make([]*simpletable.Cell, 0)
	for _, key := range keys {
		cells = append(cells, &simpletable.Cell{
			Align: simpletable.AlignCenter,
			Text:  key,
		})
	}
	t.Table.Header = &simpletable.Header{
		Cells: cells,
	}
}

func (t *Table) AddFooter(keys []string) {
	cells := make([]*simpletable.Cell, 0)
	for _, key := range keys {
		cells = append(cells, &simpletable.Cell{
			Align: simpletable.AlignCenter,
			Text:  key,
		})
	}
	t.Table.Footer = &simpletable.Footer{
		Cells: cells,
	}
}

func (t *Table) addData(data [][]interface{}) {

}

// CreateBasicFileTable creates the generic file table
func (t *Table) CreateBasicFileTable(lmap languageMap, totalFiles int) {
	data := make(TableData, 0)
	for key, val := range lmap {
		if key != "Other" {
			data = append(data, TableDatum{key, val, (float64(val) / float64(totalFiles)) * 100})
		}
	}
	t.Data = data
	t.Table = simpletable.New()
	t.AddHeader([]string{
		"FILETYPE",
		"# OF FILES",
		"% OF TOTAL",
	})
	t.AddFooter([]string{
		"", fmt.Sprintf("%d", totalFiles), "100%",
	})
	t.SetCellsForData(nil)
}
