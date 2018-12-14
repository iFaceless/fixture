package exporters

import (
	"fmt"
	"strings"
)

var sqlTemplate = `INSERT INTO %s (%s)
VALUES 
%s;
`

type SQLExporter struct{}

func NewSQLExporter() *SQLExporter {
	return new(SQLExporter)
}

func (exporter *SQLExporter) Export(tableName string, columns []string, rawRows [][][]byte) ([]byte, error) {
	rows := make([]string, 0)
	for _, rawRow := range rawRows {
		row := make([]string, len(rawRow))
		for i, colValue := range rawRow {
			if colValue == nil {
				row[i] = "NULL"
			} else {
				str := "'" + string(colValue) + "'"
				row[i] = strings.Replace(str, "\n", "\\n", len(str))
			}
		}
		rows = append(rows, "    ("+strings.Join(row, ", ")+")")
	}
	sqlExp := fmt.Sprintf(sqlTemplate, tableName, strings.Join(columns, ", "), strings.Join(rows, ",\n"))
	return []byte(sqlExp), nil
}
