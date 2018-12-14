package loaders

import (
	"encoding/json"
	"fmt"
	"os"
)

type JsonLoader struct{}

func NewJsonLoader() *JsonLoader {
	return new(JsonLoader)
}

func (loader *JsonLoader) Load(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}

	var content LoadContent
	err = json.NewDecoder(file).Decode(&content)
	if err != nil {
		panic(fmt.Sprintf("failed to load file '%s': %s", filename, err))
	}

	return genSQL(&content), nil
}
