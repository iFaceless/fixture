package loaders

import "io/ioutil"

type SQLLoader struct{}

func NewSQLLoader() *SQLLoader {
	return &SQLLoader{}
}

func (loader *SQLLoader) Load(filename string) (string, error) {
	buf, err := ioutil.ReadFile(filename)
	return string(buf), err
}
