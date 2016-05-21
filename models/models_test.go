package models

import (
	//	"time"
	//	"fmt"
	"testing"
)

func TestInitDatabase(t *testing.T) {
	err := InitDatabase("sqlite3", ":memory:", true)
	if err != nil {
		t.Error(err)
	}
}

func TestCreateUser(t *testing.T) {
	var user User
	user.Name = "test"
	err := user.SetPassword("test")
	if err != nil {
		t.Errorf("Error setting password - %s", err)
	}

	err = DB.Create(&user).Error
	if err != nil {
		t.Errorf("Error saving user profile to database - %s", err)
	}
}

func TestFetchUserAndTestPassword(t *testing.T) {
	var user User
	err := DB.Where("name = ?", "test").First(&user).Error
	if err != nil {
		t.Errorf("Error fetching test user - %s", err)
	}

	if user.CheckPassword("wrong password") {
		t.Errorf("Wrong password accepted!")
	}

	if !user.CheckPassword("test") {
		t.Errorf("Good password rejected!")
	}
}

func TestCreateMessage(t *testing.T) {

}

func TestGetMessages(t *testing.T) {

}
