package fixture

import (
	"fmt"

	"github.com/iFaceless/fixture/exporters"
)

type Exporter interface {
	Export(tableName string, columns []string, rawRows [][][]byte) ([]byte, error)
}

var exporterMap = make(map[DataFormat]Exporter)

func LookupExporter(format DataFormat) Exporter {
	exporter, ok := exporterMap[format]
	if !ok {
		panic(fmt.Sprintf("exporter not found for data format '%s'", format))
	}
	return exporter
}

func RegisterExporter(dataFmt DataFormat, ext string, exporter Exporter) {
	if oldFmt, ok := extToDataFmtMapping[ext]; ok {
		if oldFmt != dataFmt {
			panic(fmt.Sprintf("old data format of extension '%s' dose not match the new one", ext))
		}
	} else {
		extToDataFmtMapping[ext] = dataFmt
	}

	exporterMap[dataFmt] = exporter
}

func init() {
	RegisterExporter(SQL, ".sql", exporters.NewSQLExporter())
	RegisterExporter(JSON, ".json", exporters.NewJsonExporter())
	RegisterExporter(YAML, ".yml", exporters.NewYamlExporter())
}
