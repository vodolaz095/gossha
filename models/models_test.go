package models

import (
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
}

func TestFetchUserAndTestPassword(t *testing.T) {
	var user User
	err := DB.Where("name = ?", "test").First(&user).Error
	if err != nil {
		t.Errorf("Error fetching test user - %s", err)
	}
	wrong, err := user.CheckPassword("wrong password")
	if err != nil {
		t.Errorf("%s : while checking bad password", err)
	}
	if wrong != false {
		t.Errorf("wrong password accepted")
	}

	good, err := user.CheckPassword("test")
	if err != nil {
		t.Errorf("%s : while checking bad password", err)
	}
	if good != true {
		t.Errorf("good password rejected")
	}
}

func TestCreateMessage(t *testing.T) {
	t.Skipf("not implemented %s", "yet")
}

func TestGetMessages(t *testing.T) {
	t.Skipf("not implemented %s", "yet")
}
