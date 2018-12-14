package fixture

import (
	"testing"

	"path"

	"fmt"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type SuiteConfigTester struct {
	suite.Suite
	tf *TestFixture
}

func (s *SuiteConfigTester) SetupTest() {
	s.tf = &TestFixture{config: new(Config)}
}

func (s *SuiteConfigTester) Test_ConfigSchemaFilepath() {
	pth := path.Join(testDataDir, "schema.sql")
	SchemaFilepath(pth)(s.tf)
	assert.Equal(s.T(), pth, s.tf.Config().SchemaFilepath)
}

func (s *SuiteConfigTester) Test_ConfigDataDir() {
	DataDir(testDataDir)(s.tf)
	assert.Equal(s.T(), testDataDir, s.tf.Config().FixtureDataDir)
}

func (s *SuiteConfigTester) Test_ConfigDatabase() {
	rawurl := "mysql://user:password@localhost:1234/test_db?loc=Asia/Shanghai&parseTime=true"
	Database(rawurl)(s.tf)
	assert.Equal(s.T(), rawurl, s.tf.Config().DatabaseURL.String())

	rawurl = "user:password@localhost:1234/test_db?loc=Asia/Shanghai&parseTime=true"
	assert.Panics(s.T(), func() {
		Database(rawurl)(s.tf)
	})

	rawurl = "mysql://user:password@localhost:1234?loc=Asia/Shanghai&parseTime=true"
	assert.PanicsWithValue(s.T(), fmt.Sprintf("db name not found in '%s'", rawurl), func() {
		Database(rawurl)(s.tf)
	})

	rawurl = "mysql://user:password@localhost:1234/product?loc=Asia/Shanghai&parseTime=true"
	assert.PanicsWithValue(s.T(), "invalid db name 'product': test db name must starts with 'test_'", func() {
		Database(rawurl)(s.tf)
	})
}

func TestSuiteConfig(t *testing.T) {
	suite.Run(t, new(SuiteConfigTester))
}

func Test_ConfigValidate(t *testing.T) {
	conf := &Config{}

	assert.Equal(t, ErrMissingDBRawURL, conf.Validate())

	conf.DatabaseURL = new(DatabaseURL)
	conf.FixtureDataDir = "/path/to/foo"
	assert.Equal(t, ErrFixtureDataDirNotFound, conf.Validate())

	conf.FixtureDataDir = fixtureDataDir
	conf.SchemaFilepath = "/path/to/bazz.sql"
	assert.Equal(t, ErrSchemaFileNotFound, conf.Validate())

	conf.SchemaFilepath = path.Join(testDataDir, "schema.sql")
	assert.Nil(t, conf.Validate())
}
