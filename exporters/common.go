package exporters

import (
	"github.com/go-errors/errors"
	"github.com/iancoleman/orderedmap"
	"gopkg.in/yaml.v2"
)

var (
	ErrEmptyExportContent = errors.New("fixture.exporter: empty export content")
)

type SortedMap struct {
	*orderedmap.OrderedMap
}

func NewSortedMap() *SortedMap {
	return &SortedMap{orderedmap.New()}
}

func (sortedMap *SortedMap) MarshalYAML() (interface{}, error) {
	items := make([]yaml.MapItem, 0)
	for _, k := range sortedMap.Keys() {
		v, _ := sortedMap.Get(k)
		items = append(items, yaml.MapItem{Key: k, Value: v})
	}
	// Reference: https://blog.labix.org/2014/09/22/announcing-yaml-v2-for-go#mapslice
	return yaml.MapSlice(items), nil
}

type ExportContent struct {
	Table   string       `json:"table" yaml:"table"`
	Version string       `json:"version" yaml:"version"`
	Rows    []*SortedMap `json:"rows" yaml:"rows"`
}

func genExportContent(tableName string, columns []string, rawRows [][][]byte) *ExportContent {
	content := &ExportContent{
		Table:   tableName,
		Version: "1.0",
	}

	rows := make([]*SortedMap, 0)
	for _, rawRow := range rawRows {
		row := NewSortedMap()
		for i, col := range columns {
			val := rawRow[i]
			if val == nil {
				row.Set(col, nil)
			} else {
				row.Set(col, string(val))
			}
		}
		rows = append(rows, row)
	}

	content.Rows = rows

	return content
}
