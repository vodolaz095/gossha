package models

import (
	"github.com/jinzhu/gorm"
)

//DB is the handler for database
var DB *gorm.DB

// InitDatabase initialize new database connection
// the first argument is a driver to use, the second one - connection parameters
// Examples:
// InitDatabase("sqlite3","/var/lib/gossha/gossha.db")
// InitDatabase("sqlite3",":memory:")
// InitDatabase("mysql", "user:password@/dbname?charset=utf8&parseTime=True&loc=Local")
// InitDatabase("postgres", "user=gorm dbname=gorm sslmode=disable")
// InitDatabase("postgres", "postgres://pqgotest:password@localhost/pqgotest?sslmode=verify-full")
func InitDatabase(driver, dbPath string, verboseLogging bool) error {
	var db *gorm.DB
	db, err := gorm.Open(driver, dbPath)
	db.LogMode(verboseLogging)
	err = db.DB().Ping()
	if err != nil {
		return err
	}
	db.DB().SetMaxIdleConns(5)
	db.DB().SetMaxOpenConns(5)
	// Disable table name's pluralization
	db.SingularTable(true)
	err = db.AutoMigrate(&User{}, &Message{}, &Key{}, &Session{}).Error
	if err != nil {
		return err
	}
	err = db.Model(&User{}).AddUniqueIndex("idx_contact_name", "name").Error
	if err != nil {
		return err
	}
	err = db.Model(&Message{}).AddIndex("idx_message", "user_id", "created_at").Error
	if err != nil {
		return err
	}
	err = db.Model(&Key{}).AddIndex("idx_key", "user_id", "created_at").Error
	if err != nil {
		return err
	}
	err = db.Model(&Key{}).AddUniqueIndex("idx_key_content", "content", "user_id", "created_at").Error
	if err != nil {
		return err
	}
	err = db.Model(&Session{}).AddIndex("idx_session", "user_id", "created_at").Error
	if err != nil {
		return err
	}
	DB = db
	return nil
}
