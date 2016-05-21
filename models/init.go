package models

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

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
	var err error

	switch driver {
	case "sqlite3":
		db, err = gorm.Open("sqlite3", dbPath)
		if err != nil {
			return err
		}
		break
	case "mysql":
		db, err = gorm.Open("mysql", dbPath)
		if err != nil {
			return err
		}
		break
	case "postgres":
		db, err = gorm.Open("postgres", dbPath)
		if err != nil {
			return err
		}
		break
	default:
		return fmt.Errorf("Unknown database driver of %v", driver)
	}

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
