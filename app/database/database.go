package database

import (
	"fmt"

	"github.com/jinzhu/gorm"

	// Let GORM know about MySQL
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// Open creates a new connection to a MySQL database
// with the given credentials
func Open(username, password, db string) (*gorm.DB, error) {
	// Connect to the given database
	databaseHandle, err := gorm.Open("mysql", fmt.Sprintf(
		"%v:%v@/%v?charset=utf8&parseTime=True&loc=Local",
		username,
		password,
		db,
	))
	if err != nil {
		return nil, err
	}
	return databaseHandle, nil
}
