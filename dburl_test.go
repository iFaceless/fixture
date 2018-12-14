package fixture

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDatabaseURL(t *testing.T) {
	rawurl := "mysql://user:password@localhost:1234/test_db?loc=Asia/Shanghai&parseTime=true"
	parsedURL, err := Parse(rawurl)
	assert.Nil(t, err)
	assert.Equal(t, "mysql", parsedURL.Driver())

	expectedDSN := "user:password@tcp(localhost:1234)/test_db?loc=Asia%2FShanghai&parseTime=true"
	assert.Equal(t, expectedDSN, parsedURL.DSN())

	assert.Equal(t, "test_db", parsedURL.DBName())
	assert.Equal(t, rawurl, parsedURL.String())
}

func TestDatabaseURL_WithoutDBName(t *testing.T) {
	rawurl := "mysql://user:password@localhost:1234?loc=Asia/Shanghai&parseTime=true"
	parsedURL, err := Parse(rawurl)
	assert.Nil(t, err)
	assert.Equal(t, "", parsedURL.DBName())
}

func TestDatabaseURL_Failed(t *testing.T) {
	rawurl := "sqlite://user:password@localhost:1234?loc=Asia/Shanghai&parseTime=true"
	_, err := Parse(rawurl)
	assert.Equal(t, ErrDriverNotSupported, err)

	rawurl = "user:password@localhost:1234?loc=Asia/Shanghai&parseTime=true"
	_, err = Parse(rawurl)
	assert.NotNil(t, err)
}
