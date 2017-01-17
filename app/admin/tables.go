package admin

import (
	"github.com/raggaer/castro/app/database"
)

// Tables holds all the server tables
var Tables []*Table

// Table struct contains all the server tables
type Table struct {
	Database      string
	TableName     string
	Engine        string
	AutoIncrement int
	Fields        []Field
}

// Field struct contains all the fields of a table
type Field struct {
	ColumnName string
	DataType   string
}

// GetDatabaseInformation gets all the information from the database
// to populate the admin table list
func GetDatabaseInformation(database string) ([]*Table, error) {
	// Get table list
	tables, err := GetDatabaseTables(database)

	return tables, err
}

// GetFields gets all fields from a table
func (t *Table) GetFields() error {
	// Get fields
	rows, err := database.DB.Table("information_schema.columns").Select("column_name, data_type").Where("table_schema = ? AND table_name = ?", t.Database, t.TableName).Rows()

	if err != nil {
		return err
	}

	// Close rows
	defer rows.Close()

	// Loop rows
	for rows.Next() {

		// Data holder
		field := Field{}

		// Scan row
		rows.Scan(&field.ColumnName, &field.DataType)

		// Append
		t.Fields = append(t.Fields, field)
	}

	return nil
}

// GetDatabaseTables gets all the database tables
// as a map
func GetDatabaseTables(db string) ([]*Table, error) {
	// Data holder
	tableList := []*Table{}

	// Get table list
	rows, err := database.DB.Table("information_schema.tables").Select("table_name, engine, auto_increment").Where("table_schema = ?", db).Rows()

	if err != nil {
		return nil, err
	}

	// Close rows
	defer rows.Close()

	// Loop rows
	for rows.Next() {

		// Data holder
		table := &Table{
			Database: db,
			Fields:   []Field{},
		}

		// Scan row
		rows.Scan(&table.TableName, &table.Engine, &table.AutoIncrement)

		// Append table
		tableList = append(tableList, table)

		// Get table fields
		if err := table.GetFields(); err != nil {
			return nil, err
		}
	}

	return tableList, nil
}
