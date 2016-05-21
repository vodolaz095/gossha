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
	defer session.Close()

	// Set up terminal modes
	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}
	// Request pseudo terminal
	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		return *session, fmt.Errorf("request for pseudo terminal failed: %s", err)
	}

	// Start remote shell
	if err := session.Shell(); err != nil {
		return *session, fmt.Errorf("failed to start shell: %s", err)
	}
	return *session, nil
}

func TestSpawnServer(t *testing.T) {
	err := spawnServer()
	if err != nil {
		t.Error("Error spawning! -", err.Error())
	}
	t.Parallel()
	err = StartSSHD("127.0.0.1:3396")
	if err != nil {
		t.Error("Error spawning! -", err.Error())
	}
}

func TestAuthorizeViaGoodPassword(t *testing.T) {
	t.Parallel()
	time.Sleep(time.Second)
	session1, err := connect("a", "a", 3396)
	defer session1.Close()
	if err != nil {
		t.Error("Connection error:", err.Error())
	}
	session2, err := connect("b", "b", 3396)
	defer session2.Close()
	if err != nil {
		t.Error("Connection error:", err.Error())
	}
}

func TestAuthorizeViaBadPassword(t *testing.T) {
	t.Parallel()
	time.Sleep(time.Second)
	_, err := connect("a", "b", 3396)
	if err != nil {
		t.Error("gossha: We need to have error for authenticating with wrong password!")
	} else {
		t.Error(err)
	}
}

//func TestQuiteCommand(t *testing.T) {
//	session1, err := connect("a", "a", 3396)
//	defer session1.Close()

//	//todo it have to close the session
//	_, err = session1.Stdout.Write([]byte("\\q\r"))
//	//_, err = session1.Output("\\q\r")
//	if err != nil {
//		t.Error("Error sending command \\q! -", err.Error())
//	}

//	/*
//		//todo it have to return the error!
//		time.Sleep(100 * time.Millisecond)
//		_, err = session1.Stdout.Write([]byte("\\q\r"))
//		if err != nil {
//			t.Error("Error sending command \\q! -", err.Error())
//		}
//	*/
//}
