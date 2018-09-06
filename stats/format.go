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

func defaultTransform(index int, data interface{}) string {
	switch data.(type) {
	case string:
		return data.(string)
	case int:
		return fmt.Sprintf("%d", data.(int))
	case float64:
		return fmt.Sprintf("%.2f", data.(float64))
	}
	return "Bad Data"
}

func (t *Table) SetCellsForData(transform func(int, interface{}) string) {
	rows := make([][]*simpletable.Cell, 0)
	for _, val := range t.Data {
		row := make([]*simpletable.Cell, 0)
		for i, data := range val {
			if transform == nil {
				row = append(row, &simpletable.Cell{
					Align: simpletable.AlignLeft,
					Text:  defaultTransform(i, data),
				})
			} else {
				row = append(row, &simpletable.Cell{
					Align: simpletable.AlignLeft,
					Text:  transform(i, data),
				})
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
