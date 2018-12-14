package loaders

import (
	"fmt"
	"strings"
)

type LoadContent struct {
	Version string                   `yaml:"version" json:"version"`
	Table   string                   `yaml:"table" json:"table"`
	Rows    []map[string]interface{} `yaml:"rows" json:"rows"`
}

var sqlTemplate = `INSERT INTO %s (%s)
VALUES 
 %s;
`

func genSQL(content *LoadContent) string {
	if content == nil || content.Table == "" {
		return ""
	}

	if len(content.Rows) == 0 {
		return ""
	}

	columns := make([]string, 0)
	for k := range content.Rows[0] {
		columns = append(columns, k)
	}

	vals := make([]string, 0)
	for _, row := range content.Rows {
		fields := make([]string, 0)
		for _, col := range columns {
			if val, ok := row[col]; ok {
				if val == nil || val == interface{}(nil) {
					fields = append(fields, "NULL")
				} else {
					fields = append(fields, fmt.Sprintf("'%v'", val))
				}
			} else {
				panic(fmt.Sprintf("fixture.loaders: incosistent column found '%s'", col))
			}
		}

		vals = append(vals, "("+strings.Join(fields, ", ")+")")
	}

	for i, col := range columns {
		columns[i] = quote(col)
	}

	exp := fmt.Sprintf(sqlTemplate, quote(content.Table), strings.Join(columns, ", "), strings.Join(vals, ",\n"))
	return exp
}

func quote(col string) string {
	if col != "" {
		col = fmt.Sprintf("`%s`", col)
	}
	return col
}
