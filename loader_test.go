package fixture

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLookupLoader(t *testing.T) {
	_, ok := LookupLoader(SQL).(Loader)
	assert.True(t, ok)

	assert.Panics(t, func() {
		LookupLoader(100)
	})
}

const mockFooFmt DataFormat = 10

type MockLoader struct{}

func (*MockLoader) Load(filename string) (string, error) {
	return "", nil
}

func TestRegisterLoader(t *testing.T) {
	RegisterLoader(mockFooFmt, ".foo", &MockLoader{})
	_, ok := LookupLoader(mockFooFmt).(Loader)
	assert.True(t, ok)

	assert.PanicsWithValue(t, "old data format of extension '.sql' dose not match the new one", func() {
		RegisterLoader(mockFooFmt, ".sql", &MockLoader{})
	})
}
