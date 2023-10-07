package generate

import (
	"database/sql"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql" // import the MySQL driver
)

// create a global variable for the database connection
var db *sql.DB

func CreateModels(user, passsword, host, port, database string) {
	// open a connection to the database
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, passsword, host, port, database)
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	// store all models in this string
	var allModels string = "package models\n\n"

	// create a slice to hold the names of all the tables in the database
	var tableNames []string

	// query the database to get a list of all the tables
	rows, err := db.Query("SHOW TABLES")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	// iterate over the rows returned by the query,
	// adding each table name to the slice of table names
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			panic(err)
		}
		tableNames = append(tableNames, tableName)
	}

	// iterate over the slice of table names,
	// creating a struct type for each table
	for _, tableName := range tableNames {
		// query the database to get the column names and types for the table
		description := fmt.Sprintf("DESCRIBE `%s`", tableName)
		rows, err := db.Query(description)
		if err != nil {
			panic(err)
		}
		defer rows.Close()

		// create a slice to hold the column names and types
		var columns []string

		// iterate over the rows returned by the query,
		// adding each column name and type to the slice
		for rows.Next() {
			// Field Type Null Key Default Extra
			var columnName string
			var columnType string
			var columnNullable string
			var columnKey sql.NullString
			var columnDefault sql.NullString
			var columnExtra sql.NullString
			if err := rows.Scan(&columnName, &columnType,
				&columnNullable, &columnKey, &columnDefault, &columnExtra); err != nil {
				panic(err)
			}
			columns = append(columns, parseColumnType(columnName, columnType))
		}

		// create a string containing the Go code for the struct type
		var structString string
		structTypeName := parseStructTypeName(tableName)
		structString += "// " + structTypeName + " declares the model struct\n"
		structString += "type " + structTypeName + " struct {\n"
		for _, column := range columns {
			structString += "    " + column + "\n"
		}
		structString += "}\n"

		// add TableName function to override naming convention
		structString += "\t// TableName overrides the table name used by " + structTypeName + " to `" + tableName + "`\n"
		structString += "\tfunc (" + structTypeName + ") " + "TableName() string {\n"
		structString += "\t\treturn \"" + tableName + "\"\n"
		structString += "}\n\n"

		// print the struct type to the console
		fmt.Println(structString)
		// append the struct to allModels
		allModels += structString
	}

	// Write the struct to a file named "hello.txt"
	err = ioutil.WriteFile("models.go.txt", []byte(allModels), 0644)
	if err != nil {
		// log.Fatal(err)
		fmt.Println(err)
	}

}

func parseStructTypeName(tableName string) string {
	// capitalize the name of the struct type
	structTypeName := strings.Title(tableName)
	// make sure is not pluralized
	structTypeName = strings.TrimSuffix(structTypeName, "s")
	// split the tableName string on underscores
	parts := strings.Split(tableName, "_")

	// capitalize each split
	for i, part := range parts {
		parts[i] = strings.Title(part)
	}

	// join the parts back together to form the struct type name
	return strings.Join(parts, "")
}

func parseColumnName(columnName string) string {
	// capitalize the name of the struct type
	structTypeName := strings.Title(columnName)
	// make sure is not pluralized
	structTypeName = strings.TrimSuffix(structTypeName, "s")
	// split the tableName string on underscores
	parts := strings.Split(columnName, "_")

	// capitalize each split
	for i, part := range parts {
		parts[i] = strings.Title(part)
	}

	// join the parts back together to form the struct type name
	structTypeName = strings.Join(parts, "")

	// Replace "Id" with "ID" in the field name
	if strings.Contains(structTypeName, "Id") && strings.HasSuffix(structTypeName, "Id") {
		structTypeName = strings.TrimSuffix(structTypeName, "Id") + "ID"
	}

	return structTypeName
}

func parseColumnType(columnName, columnType string) string {
	// determine the Go primitive type based on the column_type
	// and the size and native MySQL type from the column_type string
	var goType string
	var size string
	var nativeType string
	startsWith := strings.Split(columnType, "(")
	switch startsWith[0] {
	case "tinyint":
		goType = "int8"
		fmt.Sscanf(columnType, "(%s)", &size)
		nativeType = "TINYINT"
	case "smallint":
		goType = "int16"
		fmt.Sscanf(columnType, "(%s)", &size)
		nativeType = "SMALLINT"
	case "mediumint":
		goType = "int32"
		fmt.Sscanf(columnType, "(%s)", &size)
		nativeType = "MEDIUMINT"
	case "int":
		goType = "int32"
		fmt.Sscanf(columnType, "(%s)", &size)
		nativeType = "INT"
	case "bigint":
		goType = "int64"
		fmt.Sscanf(columnType, "(%s)", &size)
		nativeType = "BIGINT"
	case "float":
		goType = "float32"
		fmt.Sscanf(columnType, "(%s)", &size)
		nativeType = "FLOAT"
	case "double":
		goType = "float64"
		fmt.Sscanf(columnType, "(%s)", &size)
		nativeType = "DOUBLE"
	case "decimal":
		goType = "float64"
		fmt.Sscanf(columnType, "(%s)", &size)
		nativeType = "DECIMAL"
	case "bit":
		goType = "int"
		fmt.Sscanf(columnType, "(%s)", &size)
		nativeType = "BIT"
	case "char":
		goType = "string"
		fmt.Sscanf(columnType, "(%s)", &size)
		nativeType = "CHAR"
	case "varchar":
		goType = "string"
		fmt.Sscanf(columnType, "(%s)", &size)
		nativeType = "VARCHAR"
	case "binary":
		goType = "[]byte"
		fmt.Sscanf(columnType, "(%s)", &size)
		nativeType = "BINARY"
	case "blob":
		goType = "[]byte"
		fmt.Sscanf(columnType, "(%s)", &size)
		nativeType = "BLOB"
	case "tinyblob":
		goType = "[]byte"
		fmt.Sscanf(columnType, "(%s)", &size)
		nativeType = "TINYBLOB"
	case "text":
		goType = "string"
		fmt.Sscanf(columnType, "(%s)", &size)
		nativeType = "TEXT"
	case "tinytext":
		goType = "string"
		fmt.Sscanf(columnType, "(%s)", &size)
		nativeType = "TINYTEXT"
	case "date":
		goType = "time.Time"
		fmt.Sscanf(columnType, "(%s)", &size)
		nativeType = "DATE"
	case "datetime":
		goType = "time.Time"
		fmt.Sscanf(columnType, "(%s)", &size)
		nativeType = "DATETIME"
	case "timestamp":
		goType = "time.Time"
		fmt.Sscanf(columnType, "(%s)", &size)
		nativeType = "TIMESTAMP"
	case "time":
		goType = "time.Duration"
		fmt.Sscanf(columnType, "(%s)", &size)
		nativeType = "TIMESTAMP"
	case "year":
		goType = "int"
		fmt.Sscanf(columnType, "(%s)", &size)
		nativeType = "YEAR"
	}
	// remove last ) that sometimes gets capture by the Sscanf
	size = strings.TrimSuffix(size, ")")

	var gormAnnotation string
	if len(size) > 0 {
		gormAnnotation = fmt.Sprintf("`json:\"%s\" gorm:\"column:%s;type:%s(%s)\"`", columnName, columnName, nativeType, size)
	} else {
		gormAnnotation = fmt.Sprintf("`json:\"%s\" gorm:\"column:%s;type:%s\"`", columnName, columnName, columnType)
	}

	// fix columnName to memberName
	memberName := parseColumnName(columnName)

	res := fmt.Sprintf("\t%s \t%s \t%s", memberName, goType, gormAnnotation)

	return res
}

func replaceIdWithID(fileName string) error {
	// Parse the GoLang source code
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, fileName, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	// Traverse the AST and find all the field names that end with "Id"
	ast.Inspect(f, func(n ast.Node) bool {
		// Check if the node is a field
		field, ok := n.(*ast.Field)
		if !ok {
			return true
		}

		// fmt.Println("field:", field)
		if len(field.Names) == 0 {
			return true
		}

		// Check if the field name ends with "Id"
		if !strings.HasSuffix(field.Names[0].Name, "Id") {
			return true
		}

		// Replace "Id" with "ID" in the field name
		if strings.Contains(field.Names[0].Name, "Id") && strings.HasSuffix(field.Names[0].Name, "Id") {
			field.Names[0].Name = strings.TrimSuffix(field.Names[0].Name, "Id") + "ID"
		}

		return true
	})

	// Write the modified AST back to the file
	return format.Node(os.Stdout, fset, f)
}
