package gossha

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql" //See https://github.com/jinzhu/gorm#initialize-database
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"           //See https://github.com/jinzhu/gorm#initialize-database
	_ "github.com/mattn/go-sqlite3" //See https://github.com/jinzhu/gorm#initialize-database
	"time"
)

// User represents user of chat, it is persisted in relational
// database via https://github.com/jinzhu/gorm object relational mapper
type User struct {
	ID             int64
	Name           string `sql:"size:65;unique_index"`
	Salt           string `sql:"size:65"`
	Password       string `sql:"size:65"`
	Root           bool
	LastSeenOnline time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Messages       []Message
	Keys           []Key
	Sessions       []Session
}

// SetPassword used to set password
func (u *User) SetPassword(password string) error {
	slt, err := GenSalt()
	if err != nil {
		return err
	}
	u.Salt = slt
	u.Password = Hash(fmt.Sprintf("%v%v", password, slt))
	return nil
}

// CheckPassword returns true, if we quessed it properly
func (u *User) CheckPassword(password string) bool {
	//fmt.Println("Hash    :", fmt.Sprintf("%v%v", password, u.Salt))
	//fmt.Println("Password:", u.Password)
	return u.Password == Hash(fmt.Sprintf("%v%v", password, u.Salt))
}

// IsOnline returns true, if user done any actions within 1 minute
func (u *User) IsOnline() bool {
	return time.Since(u.LastSeenOnline).Minutes() < 1
}

// CreateUser creates or updates user in database with username, password, and root permissions given
func CreateUser(name, password string, root bool) error {
	var user User
	err := DB.Table("user").FirstOrInit(&user, User{Name: name}).Error
	if err != nil {
		return err
	}
	user.Root = root
	err = user.SetPassword(password)
	if err != nil {
		return err
	}

	return DB.Table("user").Save(&user).Error
}

// BanUser removes user and all his/her messages
func BanUser(name string) error {
	var user User
	err := DB.Table("user").Where("name = ?", name).First(&user).Error
	if err != nil {
		if err == gorm.RecordNotFound {
			return fmt.Errorf("User %v not found!", name)
		}
		return err
	}
	err = DB.Delete(&user).Error
	if err != nil {
		return err
	}
	err = DB.Table("message").Where("user_id", user.ID).Delete(Message{}).Error
	if err != nil {
		return err
	}
	err = DB.Table("session").Where("user_id", user.ID).Delete(Session{}).Error
	if err != nil {
		return err
	}
	err = DB.Table("key").Where("user_id", user.ID).Delete(Key{}).Error
	if err != nil {
		return err
	}
	return nil
}

// InitDatabase initialize new database connection
// the first argument is a driver to use, the second one - connection parameters
// Examples:
// gossha.InitDatabase("sqlite3","/var/lib/gossha/gossha.db")
// gossha.InitDatabase("sqlite3",":memory:")
// gossha.InitDatabase("mysql", "user:password@/dbname?charset=utf8&parseTime=True&loc=Local")
// gossha.InitDatabase("postgres", "user=gorm dbname=gorm sslmode=disable")
// gossha.InitDatabase("postgres", "postgres://pqgotest:password@localhost/pqgotest?sslmode=verify-full")
func InitDatabase(driver, dbPath string) error {
	var db gorm.DB
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

	//logging for database
	//disabled
	db.LogMode(false)
	//enabled
	//db.LogMode(true)

	db.DB().Ping()

	db.DB().SetMaxIdleConns(5)
	db.DB().SetMaxOpenConns(5)

	// Disable table name's pluralization
	db.SingularTable(true)
	db.AutoMigrate(&User{}, &Message{}, &Key{}, &Session{})

	db.Model(&User{}).AddUniqueIndex("idx_contact_name", "name")

	db.Model(&Message{}).AddIndex("idx_message", "user_id", "created_at")

	db.Model(&Key{}).AddIndex("idx_key", "user_id", "created_at")
	db.Model(&Key{}).AddUniqueIndex("idx_key_content", "content", "user_id", "created_at")

	db.Model(&Session{}).AddIndex("idx_session", "user_id", "created_at")

	// Export globals
	Board = make(map[string]*Handler)
	DB = &db

	return nil
}

// Message represents message in chat, it is persisted in relational
// database via https://github.com/jinzhu/gorm object relational mapper
type Message struct {
	ID        int64
	User      User
	UserID    int64  // Foreign key for User (belongs to)
	Message   string `sql:"size:255"`
	IP        string `sql:"size:65"`
	Hostname  string `sql:"size:65"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Key stories the hash of public key, used by User to authorize, it is
// persisted in relational database via https://github.com/jinzhu/gorm
// object relational mapper
type Key struct {
	ID        int64
	User      User
	UserID    int64  // Foreign key for User (belongs to)
	Content   string `sql:"size:65"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Session stories Ip and Hostname, used by User to authorize, it is
// persisted in relational database via https://github.com/jinzhu/gorm
// object relational mapper
type Session struct {
	ID        int64
	User      User
	UserID    int64  // Foreign key for User (belongs to)
	IP        string `sql:"size:65"`
	Hostname  string `sql:"size:65"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Notification represents message or system notification emitted by one of Handler's,
// this struct is not persisted in database
type Notification struct {
	User             User
	Message          Message
	IsSystem         bool
	IsChat           bool
	IsPrivateMessage bool
}
