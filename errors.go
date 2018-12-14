package fixture

import "errors"

var (
	ErrFixtureDataDirNotFound = errors.New("fixture data dir not found")
	ErrSchemaFileNotFound     = errors.New("schema file not found")
	ErrMissingDBRawURL        = errors.New("database url is not configured")
	ErrDriverNotSupported     = errors.New("only mysql driver is officially supported")
)
