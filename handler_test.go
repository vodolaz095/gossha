package gossha

import (
	"net"
	"os"
	"testing"
)

type connFakeAddr struct{}

func (c connFakeAddr) Network() string {
	return ""
}
func (c connFakeAddr) String() string {
	return "127.0.0.1:22"
}

type connMetadataFake struct{}

func (c connMetadataFake) User() string {
	return "testUser"
}
func (c connMetadataFake) SessionID() []byte {
	return []byte("lalalalala")
}
func (c connMetadataFake) ClientVersion() []byte {
	return []byte("1")
}
func (c connMetadataFake) ServerVersion() []byte {
	return []byte("1")
}
func (c connMetadataFake) LocalAddr() net.Addr {
	return connFakeAddr{}
}
func (c connMetadataFake) RemoteAddr() net.Addr {
	return connFakeAddr{}
}

func TestSqlite3InitDatabase(t *testing.T) {
	err := InitDatabase("sqlite3", ":memory:")
	if err != nil {
		t.Error(err.Error())
	}
	err = CreateUser("testUser", "testPassword", false)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestSqlite3LoginByUsernameAndPassword(t *testing.T) {
	h1 := New()
	m := connMetadataFake{}
	err := h1.LoginByUsernameAndPassword(m, "testPassword")
	if err != nil {
		t.Error(err.Error())
	}
	if h1.CurrentUser.Name != "testUser" {
		t.Error("We have authorized wrong user via Handler.LoginByUsernameAndPassword!")
	}
	if h1.SessionId != "lalalalala" {
		t.Error("Session Id is not set via Handler.LoginByUsernameAndPassword!")
	}
	if h1.Ip != "127.0.0.1" {
		t.Error("Remote IP  is not set via Handler.LoginByUsernameAndPassword!")
	}

	h2 := New()
	err = h2.LoginByUsernameAndPassword(m, "wrongTestPassword")
	if err != nil {
		if err.Error() != "Wrong password for user testUser!" {
			t.Error(err.Error())
		}
	} else {
		t.Error("No error for wrong password via Handler.LoginByUsernameAndPassword!")
	}
	if h2.CurrentUser.Name != "" {
		t.Error("We have authorized wrong user via Handler.LoginByUsernameAndPassword!")
	}
}

func TestMysqlInitDatabase(t *testing.T) {
	if os.Getenv("IS_TRAVIS") == "YES" {
		err := InitDatabase("mysql", "travis:@localhost/gossha_test?charset=utf8")
		if err != nil {
			t.Error(err.Error())
		}
		err = CreateUser("testUser", "testPassword", false)
		if err != nil {
			t.Error(err.Error())
		}
	}
}

func TestMysqlLoginByUsernameAndPassword(t *testing.T) {
	if os.Getenv("IS_TRAVIS") == "YES" {
		h1 := New()
		m := connMetadataFake{}
		err := h1.LoginByUsernameAndPassword(m, "testPassword")
		if err != nil {
			t.Error(err.Error())
		}
		if h1.CurrentUser.Name != "testUser" {
			t.Error("We have authorized wrong user via Handler.LoginByUsernameAndPassword!")
		}
		if h1.SessionId != "lalalalala" {
			t.Error("Session Id is not set via Handler.LoginByUsernameAndPassword!")
		}
		if h1.Ip != "127.0.0.1" {
			t.Error("Remote IP  is not set via Handler.LoginByUsernameAndPassword!")
		}

		h2 := New()
		err = h2.LoginByUsernameAndPassword(m, "wrongTestPassword")
		if err != nil {
			if err.Error() != "Wrong password for user testUser!" {
				t.Error(err.Error())
			}
		} else {
			t.Error("No error for wrong password via Handler.LoginByUsernameAndPassword!")
		}
		if h2.CurrentUser.Name != "" {
			t.Error("We have authorized wrong user via Handler.LoginByUsernameAndPassword!")
		}
	}
}

func TestPgInitDatabase(t *testing.T) {
	if os.Getenv("IS_TRAVIS") == "YES" {
		err := InitDatabase("postgres", "postgres://postgres:@localhost/gossha_test")
		if err != nil {
			t.Error(err.Error())
		}
		err = CreateUser("testUser", "testPassword", false)
		if err != nil {
			t.Error(err.Error())
		}
	}
}

func TestPgLoginByUsernameAndPassword(t *testing.T) {
	if os.Getenv("IS_TRAVIS") == "YES" {
		h1 := New()
		m := connMetadataFake{}
		err := h1.LoginByUsernameAndPassword(m, "testPassword")
		if err != nil {
			t.Error(err.Error())
		}
		if h1.CurrentUser.Name != "testUser" {
			t.Error("We have authorized wrong user via Handler.LoginByUsernameAndPassword!")
		}
		if h1.SessionId != "lalalalala" {
			t.Error("Session Id is not set via Handler.LoginByUsernameAndPassword!")
		}
		if h1.Ip != "127.0.0.1" {
			t.Error("Remote IP  is not set via Handler.LoginByUsernameAndPassword!")
		}

		h2 := New()
		err = h2.LoginByUsernameAndPassword(m, "wrongTestPassword")
		if err != nil {
			if err.Error() != "Wrong password for user testUser!" {
				t.Error(err.Error())
			}
		} else {
			t.Error("No error for wrong password via Handler.LoginByUsernameAndPassword!")
		}
		if h2.CurrentUser.Name != "" {
			t.Error("We have authorized wrong user via Handler.LoginByUsernameAndPassword!")
		}
	}
}
