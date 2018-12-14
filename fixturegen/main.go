package main

import (
	"database/sql"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/go-errors/errors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/iFaceless/fixture"
)

var (
	dbRawURlFlag  = flag.String("url", "", "database connection url")
	tableNameFlag = flag.String("t", "", "table to be exported")
	queryFlag     = flag.String("q", "", "custom query sql")
	outputDirFlag = flag.String("o", ".", "output directory")
	extFlag       = flag.String("ext", ".yml", "output file extension (e.g. '.yml', '.json'）")
)

func main() {
	flag.Parse()
	arg := getExportArg()
	columns, rawResults := fetchQueryResults(arg)
	exportResults(arg, columns, rawResults)
	fmt.Printf("fixture.exporter: succeeded to export query results to '%s'\n", path.Join(arg.outputDir, arg.tableName+arg.ext))
}

type exportArg struct {
	dburl     *fixture.DatabaseURL
	query     string
	tableName string
	outputDir string
	ext       string
}

func getExportArg() exportArg {
	dburl, err := fixture.Parse(*dbRawURlFlag)
	exitOnError(err)

	tableName := *tableNameFlag
	if tableName == "" {
		exitOnError(errors.New("table name not specified"))
	}

	query := getValidQuery(*queryFlag, tableName)

	outputDir := *outputDirFlag
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		exitOnError(err)
	}

	ext := *extFlag
	if _, ok := fixture.LookupDataFormatByExt(ext); !ok {
		exitOnError(errors.Errorf("unsupported file extension '%s'", ext))
	}

	return exportArg{
		dburl:     dburl,
		query:     query,
		tableName: tableName,
		outputDir: outputDir,
		ext:       ext,
	}
}

func getValidQuery(query string, tableName string) string {
	if query == "" {
		query = fmt.Sprintf("SELECT * FROM %s ORDER BY id LIMIT 10", tableName)
	}
	if !strings.HasPrefix(strings.ToLower(query), "select") {
		exitOnError(errors.New("query must starts with SELECT"))
	}
	if !strings.Contains(query, tableName) {
		exitOnError(errors.New("table name not found in query sql"))
	}

	// TODO: Check complicated SQL expressions
	return query
}

func exitOnError(err error) {
	if err != nil {
		fmt.Println(err)
		flag.Usage()
		os.Exit(1)
	}
}

func fetchQueryResults(arg exportArg) ([]string, [][][]byte) {
	db, err := sql.Open(arg.dburl.Driver(), arg.dburl.DSN())
	exitOnError(err)
	defer db.Close()

	rows, err := db.Query(arg.query)
	exitOnError(err)
	defer rows.Close()

	columns, err := rows.Columns()
	exitOnError(err)

	rowValues := make([][][]byte, 0)

	for rows.Next() {
		columnValues := make([][]byte, len(columns))
		dest := make([]interface{}, len(columns))
		for i := range columnValues {
			dest[i] = &columnValues[i]
		}

		err := rows.Scan(dest...)
		exitOnError(err)

		rowValues = append(rowValues, columnValues)
	}
	return columns, rowValues
}

func exportResults(arg exportArg, cols []string, rawResults [][][]byte) {
	if len(cols) == 0 || len(rawResults) == 0 {
		return
	}

	output := getExportedContent(arg, cols, rawResults)
	if len(output) == 0 {
		fmt.Println("export failed: empty ouput from exporter")
		os.Exit(1)
	}

	err := ioutil.WriteFile(
		path.Join(arg.outputDir, arg.tableName+arg.ext),
		output,
		os.ModePerm,
	)
	exitOnError(err)
}

func getExportedContent(arg exportArg, cols []string, rawResults [][][]byte) []byte {
	dataFmt, _ := fixture.LookupDataFormatByExt(arg.ext)
	exporter := fixture.LookupExporter(dataFmt)
	output, err := exporter.Export(arg.tableName, cols, rawResults)
	exitOnError(err)
	return output
}
