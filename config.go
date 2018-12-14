package fixture

import (
	"fmt"

	"strings"

	"log"
)

type Config struct {
	DatabaseURL    *DatabaseURL
	FixtureDataDir string
	SchemaFilepath string
}

func (c *Config) Validate() error {
	if c.DatabaseURL == nil {
		return ErrMissingDBRawURL
	}

	if !isPathExist(c.FixtureDataDir) {
		return ErrFixtureDataDirNotFound
	}

	if !isPathExist(c.SchemaFilepath) {
		return ErrSchemaFileNotFound
	}

	return nil
}

type Option func(*TestFixture)

func SchemaFilepath(p string) Option {
	return func(tf *TestFixture) {
		tf.config.SchemaFilepath = p
	}
}

func DataDir(dir string) Option {
	return func(tf *TestFixture) {
		tf.config.FixtureDataDir = dir
	}
}

func Database(rawurl string) Option {
	return func(tf *TestFixture) {
		url, err := Parse(rawurl)
		panicOnErr(err)

		// parse test db name from url
		dbName := url.DBName()
		if dbName == "" {
			panic(fmt.Sprintf("db name not found in '%s'", rawurl))
		}

		if !strings.HasPrefix(dbName, "test_") {
			panic(fmt.Sprintf("invalid db name '%s': test db name must starts with 'test_'", dbName))
		}

		log.Printf("config.Database: database config is '%s'", url)
		tf.config.DatabaseURL = url
	}
}
