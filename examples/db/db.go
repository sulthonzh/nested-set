package db

import (
	"errors"
	"fmt"
	"os"

	"gitlab.com/sulthonzh/scraperpath-nested-set/examples/config"

	"github.com/jinzhu/gorm"
	// spesial import ::
	_ "github.com/jinzhu/gorm/dialects/mysql"
	// _ "github.com/jinzhu/gorm/dialects/postgres"
	// _ "github.com/jinzhu/gorm/dialects/sqlite"
)

// InitDB ::
func InitDB() (db *gorm.DB, err error) {
	dbConfig := config.Config.DB
	if dbConfig.Adapter == "mysql" {
		db, err = gorm.Open("mysql", fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=True&loc=Local", dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Name))
		// DB = DB.Set("gorm:table_options", "CHARSET=utf8")
	} else if dbConfig.Adapter == "postgres" {
		db, err = gorm.Open("postgres", fmt.Sprintf("postgres://%v:%v@%v/%v?sslmode=disable", dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Name))
	} else if dbConfig.Adapter == "sqlite" {
		db, err = gorm.Open("sqlite3", fmt.Sprintf("%v/%v", os.TempDir(), dbConfig.Name))
	} else {
		err = errors.New("not supported database adapter")
	}

	return
}
