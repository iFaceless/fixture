package fixture

import (
	"database/sql"
	"fmt"
	"log"
	"path"
	"strings"
)

const (
	defaultSchemaName = "schema.sql"
)

type DataFormat int

func (f DataFormat) String() string {
	var ext string
	for k, v := range extToDataFmtMapping {
		if v == f {
			ext = k
		}
	}

	return strings.TrimLeft(ext, ".")
}

const (
	SQL DataFormat = iota
	YAML
	JSON
)

var (
	// extToDataFmtMapping maps extension to data format
	extToDataFmtMapping = map[string]DataFormat{}
)

func LookupDataFormatByExt(ext string) (DataFormat, bool) {
	ret, ok := extToDataFmtMapping[ext]
	return ret, ok
}

type TestFixture struct {
	config *Config
	tables []*table
}

func New(opts ...Option) *TestFixture {
	defaultConfig := &Config{
		FixtureDataDir: ".",
		SchemaFilepath: path.Join(".", defaultSchemaName),
	}
	tf := &TestFixture{
		config: defaultConfig,
	}

	for _, opt := range opts {
		opt(tf)
	}

	panicOnErr(tf.config.Validate())

	tf.tables = parseSchemaFile(tf.config.SchemaFilepath)
	tf.createTables()

	return tf
}

func (tf *TestFixture) String() string {
	return fmt.Sprintf("TestFixture(url='%s')", tf.config.DatabaseURL)
}

func (tf *TestFixture) Config() *Config {
	return tf.config
}

func (tf *TestFixture) TableNames() []string {
	names := make([]string, 0)
	for _, tb := range tf.tables {
		names = append(names, tb.name)
	}
	return names
}

func (tf *TestFixture) Use(tableNames ...string) *Scope {
	selectedTables := make([]*table, 0)
	for _, name := range tableNames {
		table := tf.lookupTable(name)
		if table == nil {
			panic(fmt.Sprintf("table '%s' not found", name))
		}

		selectedTables = append(selectedTables, table)
	}

	return newScope(tf, selectedTables)
}

// DropTables drops all the test tables
func (tf *TestFixture) DropTables() {
	log.Printf("fixture: drop %d tables", len(tf.tables))
	db := getDB(tf)
	defer db.Close()

	for _, tb := range tf.tables {
		_, err := db.Exec("DROP TABLE " + tb.name)
		if err != nil {
			log.Printf("fixture: failed to drop table '%s': %s", tb.name, err)
		}
	}
}

func (tf *TestFixture) createTables() {
	log.Printf("fixture: create %d tables", len(tf.tables))
	db := getDB(tf)
	defer db.Close()

	tx, err := db.Begin()
	panicOnErr(err)

	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
			panic(err)
		}
	}()

	for _, tb := range tf.tables {
		_, err := tx.Exec(tb.createSQL)
		if err != nil {
			if strings.Contains(fmt.Sprintf("%s", err), "already exists") {
				log.Printf("fixture: table '%s' already existed, try to clear existing data now", tb.name)
				db.Exec("TRUNCATE TABLE " + tb.name)
			}
		}
	}

	tx.Commit()
}

func (tf *TestFixture) lookupTable(name string) *table {
	for _, item := range tf.tables {
		if item.name == name {
			return item
		}
	}
	return nil
}

type Scope struct {
	tf             *TestFixture
	selectedTables []*table
}

func newScope(tf *TestFixture, tables []*table) *Scope {
	scope := &Scope{tf, tables}
	scope.insertFixtureData()
	return scope
}

func (s *Scope) Test(testFunc func()) {
	defer s.Clear()
	testFunc()
}

// Clear just drop the selected tables, simple and clear
func (s *Scope) Clear() {
	log.Printf("fixture: clear %d selected tables", len(s.selectedTables))
	db := getDB(s.tf)
	defer db.Close()

	for _, tb := range s.selectedTables {
		_, err := db.Exec("TRUNCATE TABLE " + tb.name)
		if err != nil {
			log.Printf("fixture: failed to clear table '%s': %s", tb.name, err)
		}
	}
}

func (s *Scope) insertFixtureData() {
	db := getDB(s.tf)
	defer db.Close()

	tx, err := db.Begin()
	panicOnErr(err)

	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
			panic(err)
		}
	}()

	for _, tb := range s.selectedTables {
		fixtureData := findFixtureData(s.tf.config.FixtureDataDir, tb)
		if fixtureData == nil {
			log.Printf("failed to find fixture data for table '%s'", tb.name)
			continue
		}

		log.Printf("insert fixture data for table '%s' from file '%s'", tb.name, fixtureData.Path)
		loader := LookupLoader(fixtureData.Format)
		sqlStr, err := loader.Load(fixtureData.Path)
		panicOnErr(err)

		if sqlStr != "" {
			_, err = tx.Exec(sqlStr)
			if err != nil {
				log.Panicf("failed to insert fixture data for table '%s': %s", tb.name, err)
			}
		}
	}

	tx.Commit()
}

func getDB(tf *TestFixture) *sql.DB {
	db, err := sql.Open(tf.config.DatabaseURL.Driver(), tf.config.DatabaseURL.DSN())
	panicOnErr(err)
	return db
}

type table struct {
	name      string
	createSQL string
}

type fixtureData struct {
	Path   string
	Format DataFormat
}
