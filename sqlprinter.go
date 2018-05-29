package sqlprinter

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/juju/errors"
	"github.com/olekukonko/tablewriter"
)

// Table prints the sql rows in a tabular structure
func Table(result *sql.Rows) error {
	var err error
	cols, err := result.Columns()
	if err != nil {
		return errors.Annotate(err, "could not get database result columns")
	}
	headers := []string{}
	rowContainer := make([]interface{}, 0)
	for _, col := range cols {
		headers = append(headers, col)
		col := new([]byte)
		rowContainer = append(rowContainer, col)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(headers)
	table.SetAutoFormatHeaders(false)

	for result.Next() {
		rowData := []string{}
		err = result.Scan(rowContainer...)
		if err != nil {
			return errors.Annotate(err, "failed parsing result row")
		}
		for _, elem := range rowContainer {
			if elem == nil {
				rowData = append(rowData, "<NULL>")
				continue
			}
			if val, ok := elem.(*[]uint8); ok {
				rowData = append(rowData, fmt.Sprintf("%s", string(*val)))
				continue
			}
		}
		table.Append(rowData)
	}

	table.Render()
	return nil
}
