package fixture

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLookupExporter(t *testing.T) {
	_, ok := LookupExporter(SQL).(Exporter)
	assert.True(t, ok)

	assert.Panics(t, func() {
		LookupLoader(100)
	})
}

const mockBarFmt DataFormat = 10

type MockExporter struct{}

func (*MockExporter) Export(tableName string, columns []string, rawRows [][][]byte) ([]byte, error) {
	return nil, nil
}

func TestRegisterExporter(t *testing.T) {
	RegisterExporter(mockBarFmt, ".foo", &MockExporter{})
	_, ok := LookupExporter(mockBarFmt).(Exporter)
	assert.True(t, ok)

	assert.PanicsWithValue(t, "old data format of extension '.sql' dose not match the new one", func() {
		RegisterExporter(mockBarFmt, ".sql", &MockExporter{})
	})
}
