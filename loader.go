package fixture

import (
	"fmt"

	"github.com/iFaceless/fixture/loaders"
)

type Loader interface {
	Load(filename string) (string, error)
}

var loaderMap = make(map[DataFormat]Loader)

func LookupLoader(dataFmt DataFormat) Loader {
	loader, ok := loaderMap[dataFmt]
	if !ok {
		panic(fmt.Sprintf("loader not found for data format '%s'", dataFmt))
	}
	return loader
}

func RegisterLoader(dataFmt DataFormat, ext string, loader Loader) {
	if oldFmt, ok := extToDataFmtMapping[ext]; ok {
		if oldFmt != dataFmt {
			panic(fmt.Sprintf("old data format of extension '%s' dose not match the new one", ext))
		}
	} else {
		extToDataFmtMapping[ext] = dataFmt
	}
	loaderMap[dataFmt] = loader
}

func init() {
	RegisterLoader(SQL, ".sql", loaders.NewSQLLoader())
	RegisterLoader(YAML, ".yml", loaders.NewYamlLoader())
	RegisterLoader(JSON, ".json", loaders.NewJsonLoader())
}
