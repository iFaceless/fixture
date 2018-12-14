package loaders

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type YamlLoader struct{}

func NewYamlLoader() *YamlLoader {
	return new(YamlLoader)
}

func (loader *YamlLoader) Load(filename string) (string, error) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}

	var content LoadContent

	err = yaml.Unmarshal(buf, &content)
	if err != nil {
		panic(fmt.Sprintf("failed to load file '%s': %s", filename, err))
	}

	return genSQL(&content), nil
}
