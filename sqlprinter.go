package sqlprinter

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/jojomi/go-spew/spew"
	"github.com/juju/errors"
	"github.com/olekukonko/tablewriter"
)

func Table(result *sql.Rows) error {
	var err error
	cols, err := result.Columns()
	if err != nil {
		return errors.Annotate(err, "could not get database result columns")
	}
	headers := []string{}
	rowContainer := make([]interface{}, 0)
	// colTypes, _ := result.ColumnTypes()
	for _, col := range cols {
		headers = append(headers, col)
		/*scanType := colTypes[i].ScanType()
		switch scanType.Kind() {
		case reflect.Uint32:
			colInt := new(uint32)
			rowContainer = append(rowContainer, colInt)
		case reflect.Slice:
			colInt := new([]byte)
			rowContainer = append(rowContainer, colInt)
		default:
			return fmt.Errorf("unknown column type: %s for column %s", scanType.Kind(), col)
		}*/
		colInt := new(interface{})
		rowContainer = append(rowContainer, colInt)
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
			spew.Dump(elem)
			/*if val, ok := elem.(*uint32); ok {
				rowData = append(rowData, fmt.Sprintf("%d", *val))
				continue
			}*/
			if val, ok := elem.(*[]uint8); ok {
				rowData = append(rowData, fmt.Sprintf("%s", string(*val)))
				continue
			}
			if val, ok := elem.(*time.Time); ok {
				rowData = append(rowData, fmt.Sprintf("%s", (*val).Format("2006")))
				continue
			}

			/// rowData = append(rowData, fmt.Sprintf("%+v", elem))
		}
		table.Append(rowData)
	}

	table.Render()
	return nil
}
