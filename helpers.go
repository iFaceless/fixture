package fixture

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"
)

func isPathExist(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

func panicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}

func findFixtureData(fixtureDataDir string, table *table) *fixtureData {
	if !isPathExist(fixtureDataDir) {
		panic(ErrFixtureDataDirNotFound)
	}

	possibleNames := make([]string, 0)
	for ext := range extToDataFmtMapping {
		possibleNames = append(possibleNames, table.name+ext)
	}

	var found string
	for _, name := range possibleNames {
		absPath := path.Join(fixtureDataDir, name)
		if isPathExist(absPath) {
			if found != "" {
				panic(fmt.Sprintf("multiple formats of fixture data found for table '%s'", table.name))
			}

			found = absPath
		}
	}

	if found == "" {
		panic(fmt.Sprintf("fixture data not found for table '%s'", table.name))
	}

	return &fixtureData{
		Format: extToDataFmtMapping[path.Ext(found)],
		Path:   found,
	}
}

var rule = regexp.MustCompile("CREATE\\s.*TABLE\\s(.*)\\(.*")

func parseSchemaFile(filename string) []*table {
	buf, err := ioutil.ReadFile(filename)
	panicOnErr(err)

	tables := make([]*table, 0)
	for _, createSQL := range strings.Split(string(buf), ";") {
		createSQL = strings.TrimSpace(createSQL)
		if createSQL == "" {
			continue
		}

		groups := rule.FindStringSubmatch(createSQL)
		if len(groups) != 2 {
			panic(fmt.Sprintf("cannot extract table name from sql '%s'", createSQL))
		}

		tb := &table{
			name:      trimTableName(groups[1]),
			createSQL: createSQL,
		}

		tables = append(tables, tb)
	}

	return tables
}

func trimTableName(n string) string {
	n = strings.Replace(n, "`", "", len(n))
	n = strings.Replace(n, "'", "", len(n))
	n = strings.Replace(n, "\"", "", len(n))
	return strings.TrimSpace(n)
}
