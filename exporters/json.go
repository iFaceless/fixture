package exporters

import "encoding/json"

type JsonExporter struct{}

func NewJsonExporter() *JsonExporter {
	return new(JsonExporter)
}

func (exporter *JsonExporter) Export(tableName string, columns []string, rawRows [][][]byte) ([]byte, error) {
	content := genExportContent(tableName, columns, rawRows)
	if content == nil {
		return nil, ErrEmptyExportContent
	}

	output, err := json.MarshalIndent(content, "", "  ")
	if err != nil {
		return nil, err
	}

	return output, nil
}
