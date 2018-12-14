package exporters

import "gopkg.in/yaml.v2"

type YamlExporter struct{}

func NewYamlExporter() *YamlExporter {
	return new(YamlExporter)
}

func (exporter *YamlExporter) Export(tableName string, columns []string, rawRows [][][]byte) ([]byte, error) {
	content := genExportContent(tableName, columns, rawRows)
	if content == nil {
		return nil, ErrEmptyExportContent
	}

	output, err := yaml.Marshal(content)
	if err != nil {
		return nil, err
	}

	return output, nil
}
