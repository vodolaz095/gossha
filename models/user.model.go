package models

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
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
		if err == gorm.ErrRecordNotFound {
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
