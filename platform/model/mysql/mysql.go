package mysql

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/obgnail/plugin-platform/platform/config"
	"github.com/obgnail/plugin-platform/utils/errors"
	"time"
)

var (
	DB *gorm.DB

	RecordNotFound = gorm.ErrRecordNotFound
)

func InitDB() (err error) {
	addr := config.StringOrPanic("mysql_address")
	user := config.StringOrPanic("mysql_user")
	password := config.StringOrPanic("mysql_password")
	dbName := config.StringOrPanic("mysql_db_name")
	DB, err = buildDBM(addr, user, password, dbName)
	if err != nil {
		return errors.Trace(err)
	}
	return nil
}

func buildDBM(address string, user string, password string, dbName string) (*gorm.DB, error) {
	var err error
	var db *gorm.DB
	str := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, address, dbName)
	db, err = gorm.Open("mysql", str)
	if err != nil {
		time.Sleep(time.Second * 5)
		panic(err)
	}
	maxIdle := config.IntOrPanic("mysql_db_max_idle")
	maxOpen := config.IntOrPanic("mysql_db_max_open")
	db.DB().SetMaxIdleConns(maxIdle)
	db.DB().SetMaxOpenConns(maxOpen)
	db.DB().SetConnMaxLifetime(time.Hour)
	db.SingularTable(true)
	db.LogMode(false)

	return db, nil
}

func Transaction(callback func(db *gorm.DB) error) error {
	if err := DB.Transaction(callback); err != nil {
		return errors.Trace(err)
	}
	return nil
}
