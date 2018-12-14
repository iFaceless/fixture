package fixture

import (
	"path"
	"testing"

	"runtime"

	"github.com/stretchr/testify/assert"
)

var (
	testDataDir    = path.Join(curDir(), "testdata")
	fixtureDataDir = path.Join(testDataDir, "fixtures")
)

func curDir() string {
	_, f, _, _ := runtime.Caller(1)
	curDir := path.Dir(f)
	return curDir
}

func Test_parseSchemaFile(t *testing.T) {
	ret := parseSchemaFile(path.Join(testDataDir, "schema.sql"))
	assert.Equal(t, 4, len(ret))
	assert.Equal(t, "user", ret[0].name)
	assert.Equal(t, "task", ret[1].name)
	assert.Equal(t, "foo", ret[2].name)
}

func Test_findFixtureData_YAML_OK(t *testing.T) {
	targetTable := &table{name: "user"}
	expectedResult := &fixtureData{
		Path:   path.Join(fixtureDataDir, "user.yml"),
		Format: YAML,
	}
	ret := findFixtureData(fixtureDataDir, targetTable)
	assert.Equal(t, expectedResult, ret)
}

func Test_findFixtureData_JSON_OK(t *testing.T) {
	targetTable := &table{name: "foo"}
	expectedResult := &fixtureData{
		Path:   path.Join(fixtureDataDir, "foo.json"),
		Format: JSON,
	}
	ret := findFixtureData(fixtureDataDir, targetTable)
	assert.Equal(t, expectedResult, ret)
}

func Test_findFixtureData_SQL_OK(t *testing.T) {
	targetTable := &table{name: "bar"}
	expectedResult := &fixtureData{
		Path:   path.Join(fixtureDataDir, "bar.sql"),
		Format: SQL,
	}
	ret := findFixtureData(fixtureDataDir, targetTable)
	assert.Equal(t, expectedResult, ret)
}

func Test_findFixtureData_Panics(t *testing.T) {
	assert.PanicsWithValue(t, "multiple formats of fixture data found for table 'beep'", func() {
		findFixtureData(fixtureDataDir, &table{name: "beep"})
	})

	assert.PanicsWithValue(t, "fixture data not found for table 'hidden'", func() {
		findFixtureData(fixtureDataDir, &table{name: "hidden"})
	})
}

func Test_isPathExist(t *testing.T) {
	assert.True(t, isPathExist(fixtureDataDir))
	assert.False(t, isPathExist(fixtureDataDir+"/file-not-found.txt"))
}
