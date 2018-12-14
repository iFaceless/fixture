package fixture

import (
	"database/sql"
	"fmt"
	"os"
	"path"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var schemaFilepath = path.Join(testDataDir, "schema.sql")

type SuiteTestFixtureTester struct {
	suite.Suite
	tf *TestFixture
	db *sql.DB
}

func (s *SuiteTestFixtureTester) SetupSuite() {
	s.tf = New(
		SchemaFilepath(schemaFilepath),
		DataDir(fixtureDataDir),
		Database(getDBRawURL()),
	)
	s.db = openDB(s.tf)

	for _, name := range s.tf.TableNames() {
		assert.True(s.T(), isTableExistInDB(s.db, name))
	}
}

func (s *SuiteTestFixtureTester) TearDownSuite() {
	s.tf.DropTables()
	for _, nm := range s.tf.TableNames() {
		assert.False(s.T(), isTableExistInDB(s.db, nm))
	}

	s.db.Close()
}

func (s *SuiteTestFixtureTester) TestGetString() {
	assert.NotNil(s.T(), s.tf)
	assert.Equal(s.T(), fmt.Sprintf("TestFixture(url='%s')", getDBRawURL()), s.tf.String())
}

func (s *SuiteTestFixtureTester) TestGetConfig() {
	conf := s.tf.Config()
	assert.NotNil(s.T(), conf)
	assert.Equal(s.T(), schemaFilepath, conf.SchemaFilepath)
	assert.Equal(s.T(), fixtureDataDir, conf.FixtureDataDir)
	assert.Equal(s.T(), getDBRawURL(), conf.DatabaseURL.String())
}

func (s *SuiteTestFixtureTester) TestUse_TableNotFound() {
	assert.PanicsWithValue(s.T(), "table 'missing_table' not found", func() {
		s.tf.Use("missing_table")
	})
}

func (s *SuiteTestFixtureTester) TestUse_TableWithSQLFixtureData() {
	targetTable := "bar"

	// Make sure table is empty
	assert.True(s.T(), isTableExistInDB(s.db, targetTable))
	assert.Equal(s.T(), 0, countTable(s.db, targetTable))

	scope := s.tf.Use(targetTable)
	assert.NotNil(s.T(), scope)
	// After data inserted
	assert.Equal(s.T(), 2, countTable(s.db, targetTable))

	scope.Clear()
	// Here, table is clean
	assert.Equal(s.T(), 0, countTable(s.db, targetTable))
}

func (s *SuiteTestFixtureTester) TestUse_TableWithJSONFixtureData() {
	targetTable := "foo"

	assert.True(s.T(), isTableExistInDB(s.db, targetTable))
	assert.Equal(s.T(), 0, countTable(s.db, targetTable))

	scope := s.tf.Use(targetTable)
	assert.NotNil(s.T(), scope)
	assert.Equal(s.T(), 2, countTable(s.db, targetTable))

	scope.Clear()
	assert.Equal(s.T(), 0, countTable(s.db, targetTable))
}

func (s *SuiteTestFixtureTester) TestUse_TableWithYAMLFixtureData() {
	targetTable := "user"

	assert.True(s.T(), isTableExistInDB(s.db, targetTable))
	assert.Equal(s.T(), 0, countTable(s.db, targetTable))

	scope := s.tf.Use(targetTable)
	assert.NotNil(s.T(), scope)
	assert.Equal(s.T(), 2, countTable(s.db, targetTable))

	scope.Clear()
	assert.Equal(s.T(), 0, countTable(s.db, targetTable))
}

func (s *SuiteTestFixtureTester) TestUse_WithAutoClear() {
	targetTable := "user"

	assert.True(s.T(), isTableExistInDB(s.db, targetTable))
	assert.Equal(s.T(), 0, countTable(s.db, targetTable))

	s.tf.Use(targetTable).Test(func() {
		assert.Equal(s.T(), 2, countTable(s.db, targetTable))
	})

	assert.Equal(s.T(), 0, countTable(s.db, targetTable))
}

func TestSuiteTestFixture(t *testing.T) {
	suite.Run(t, new(SuiteTestFixtureTester))
}

func getDBRawURL() string {
	dsnFmt := "mysql://%s:%s@%s/%s?charset=utf8&parseTime=true&loc=Asia/Shanghai"
	return fmt.Sprintf(dsnFmt,
		os.Getenv("MYSQL_USERNAME"),
		os.Getenv("MYSQL_PASSWORD"),
		os.Getenv("MYSQL_HOST"),
		os.Getenv("MYSQL_DATABASE"),
	)
}

func openDB(tf *TestFixture) *sql.DB {
	db, _ := sql.Open(tf.Config().DatabaseURL.Driver(), tf.Config().DatabaseURL.DSN())
	return db
}

func isTableExistInDB(db *sql.DB, name string) bool {
	_, err := db.Exec(fmt.Sprintf("SELECT 1 FROM %s LIMIT 1", name))
	return err == nil
}

func countTable(db *sql.DB, name string) int {
	row := db.QueryRow(fmt.Sprintf("SELECT COUNT(1) FROM %s", name))

	var count int
	row.Scan(&count)
	return count
}
