package ssh

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/vodolaz095/gossha/config"
	"github.com/vodolaz095/gossha/models"
	"golang.org/x/crypto/ssh"
)

////http://godoc.org/golang.org/x/crypto/ssh#Client

func spawnServer() error {
	config.RuntimeConfig = new(config.Config)
	config.RuntimeConfig.Port = 3396
	config.RuntimeConfig.Debug = false
	config.RuntimeConfig.Driver = "sqlite3"
	config.RuntimeConfig.ConnectionString = ":memory:"
	if os.Getenv("IS_TRAVIS") == "YES" {
		config.RuntimeConfig.SSHPublicKeyPath = "/home/travis/gopath/src/github.com/vodolaz095/gossha/test/gossha_test.pub"
		config.RuntimeConfig.SSHPrivateKeyPath = "/home/travis/gopath/src/github.com/vodolaz095/gossha/test/gossha_test"
	} else {
		config.RuntimeConfig.SSHPublicKeyPath, _ = config.GetPublicKeyPath()
		config.RuntimeConfig.SSHPrivateKeyPath, _ = config.GetPrivateKeyPath()
	}

	config.RuntimeConfig.Homedir, _ = config.GetHomeDir()
	config.RuntimeConfig.ExecuteOnMessage = ""
	config.RuntimeConfig.ExecuteOnPrivateMessage = ""
	err := models.InitDatabase("sqlite3", ":memory:", true)
	if err != nil {
		return err
	}
	err = models.CreateUser("a", "a", false)
	if err != nil {
		return err
	}
	err = models.CreateUser("b", "b", false)
	if err != nil {
		return err
	}

	return nil
}

func connect(username, password string, port int) (ssh.Session, error) {
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
	}
	client, err := ssh.Dial("tcp", fmt.Sprintf("localhost:%v", port), config)
	if err != nil {
		return ssh.Session{}, err
	}

	session, err := client.NewSession()
	if err != nil {
		return ssh.Session{}, err
	}
	return *session, err
}

func TestSpawnServer(t *testing.T) {
	err := spawnServer()
	if err != nil {
		t.Error("Error spawning! -", err.Error())
	}
	t.Parallel()
	go func() {
		err = StartSSHD("127.0.0.1:3396")
		if err != nil {
			t.Error("Error spawning! -", err.Error())
		}
	}()
}

func TestAuthorizeViaGoodPasswordForUser1(t *testing.T) {
	t.Parallel()
	time.Sleep(100 * time.Millisecond)
	session, err := connect("a", "a", 3396)
	if err != nil {
		t.Errorf("Connection error %s", err)
	}
	err = session.Close()
	if err != nil {
		t.Errorf("Error closing session %s", err)
	}
}

func TestAuthorizeViaGoodPasswordForUser2(t *testing.T) {
	t.Parallel()
	time.Sleep(100 * time.Millisecond)
	session, err := connect("b", "b", 3396)
	if err != nil {
		t.Errorf("Connection error %s", err)
	}
	err = session.Close()
	if err != nil {
		t.Errorf("Error closing session %s", err)
	}
}

func TestAuthorizeViaBadPassword(t *testing.T) {
	t.Parallel()
	time.Sleep(100 * time.Millisecond)
	_, err := connect("a", "b", 3396)
	if err != nil {
		if err.Error() != "ssh: handshake failed: ssh: unable to authenticate, attempted methods [none password], no supported methods remain" {
			t.Errorf("gossha: Wrong error: %s", err)
		}
	} else {
		t.Error("Error have to be thrown!")
	}
}

func TestSendMessage(t *testing.T) {
	t.Parallel()
	time.Sleep(100 * time.Millisecond)
	session, err := connect("a", "a", 3396)
	if err != nil {
		t.Errorf("Connection error %s", err)
	}
	err = session.Start("some test message\r")
	if err != nil {
		t.Errorf("Sending message error %s", err)
	}
}
