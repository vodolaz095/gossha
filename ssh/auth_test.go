package ssh

// TODO - test started to hang

//
//import (
//	"fmt"
//	"os"
//	"testing"
//	"time"
//
//	"github.com/vodolaz095/gossha/config"
//	"github.com/vodolaz095/gossha/models"
//	"golang.org/x/crypto/ssh"
//)
//
//////http://godoc.org/golang.org/x/crypto/ssh#Client
//
//const testPort = 3391
//
//func spawnServer() error {
//	config.RuntimeConfig = new(config.Config)
//	config.RuntimeConfig.Port = testPort
//	config.RuntimeConfig.Debug = false
//	config.RuntimeConfig.Driver = "sqlite3"
//	config.RuntimeConfig.ConnectionString = ":memory:"
//	if os.Getenv("IS_TRAVIS") == "YES" {
//		config.RuntimeConfig.SSHPublicKeyPath = "/home/travis/gopath/src/github.com/vodolaz095/gossha/test/gossha_test.pub"
//		config.RuntimeConfig.SSHPrivateKeyPath = "/home/travis/gopath/src/github.com/vodolaz095/gossha/test/gossha_test"
//	} else {
//		config.RuntimeConfig.SSHPublicKeyPath, _ = config.GetPublicKeyPath()
//		config.RuntimeConfig.SSHPrivateKeyPath, _ = config.GetPrivateKeyPath()
//	}
//
//	config.RuntimeConfig.Homedir, _ = config.GetHomeDir()
//	config.RuntimeConfig.ExecuteOnMessage = ""
//	config.RuntimeConfig.ExecuteOnPrivateMessage = ""
//	err := models.InitDatabase("sqlite3", ":memory:", true)
//	if err != nil {
//		return err
//	}
//	err = models.CreateUser("a", "a", false)
//	if err != nil {
//		return err
//	}
//	err = models.CreateUser("b", "b", false)
//	if err != nil {
//		return err
//	}
//	msg := models.Message{
//		ID:        1,
//		UserID:    1,
//		Message:   "test",
//		IP:        "127.0.0.1",
//		CreatedAt: time.Now(),
//		UpdatedAt: time.Now(),
//	}
//	return models.DB.Save(&msg).Error
//}
//
//func connect(username, password string, port int) (client *ssh.Client, sess *ssh.Session, err error) {
//	config := &ssh.ClientConfig{
//		User: username,
//		Auth: []ssh.AuthMethod{
//			ssh.Password(password),
//		},
//		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
//	}
//	client, err = ssh.Dial("tcp", fmt.Sprintf("127.0.0.1:%v", port), config)
//	if err != nil {
//		return
//	}
//	//defer client.Close()
//	sess, err = client.NewSession()
//	return
//}
//
//func TestSpawnServer(t *testing.T) {
//	err := spawnServer()
//	if err != nil {
//		t.Error("Error spawning! -", err.Error())
//	}
//	t.Logf("Server prepared, starting it on 127.0.0.1:%v...", testPort)
//	t.Parallel()
//	err = StartSSHD(fmt.Sprintf("127.0.0.1:%v", testPort))
//	if err != nil {
//		t.Error("Error spawning! -", err.Error())
//	}
//}
//
//func TestAuthorizeViaGoodPasswordForUser1(t *testing.T) {
//	t.Parallel()
//	time.Sleep(100 * time.Millisecond)
//	client, session, err := connect("a", "a", testPort)
//	if err != nil {
//		t.Errorf("Connection error %s", err)
//		return
//	}
//	t.Logf("User a/a authorized!")
//	err = session.Close()
//	if err != nil {
//		t.Errorf("Error closing session %s", err)
//	}
//	t.Logf("User a/a closed session!")
//	err = client.Close()
//	if err != nil {
//		t.Errorf("%s : while closing clinet", err)
//	}
//	t.Logf("User a/a closed client")
//}
//
//func TestAuthorizeViaGoodPasswordForUser2(t *testing.T) {
//	t.Parallel()
//	time.Sleep(100 * time.Millisecond)
//	client, session, err := connect("b", "b", testPort)
//	if err != nil {
//		t.Errorf("Connection error %s", err)
//		return
//	}
//	t.Logf("User b/b authorized!")
//	err = session.Close()
//	if err != nil {
//		t.Errorf("Error closing session %s", err)
//	}
//	t.Logf("User b/b closed session!")
//	err = client.Close()
//	if err != nil {
//		t.Errorf("%s : while closing clinet", err)
//	}
//	t.Logf("User a/a closed client")
//
//}
//
//func TestAuthorizeViaBadPassword(t *testing.T) {
//	//t.Parallel()
//	time.Sleep(100 * time.Millisecond)
//	_, _, err := connect("a", "b", testPort)
//	if err != nil {
//		t.Logf("%s : error reported when trying to authorize with bad password", err)
//		if err.Error() != "ssh: handshake failed: ssh: unable to authenticate, attempted methods [none password], no supported methods remain" || err.Error() != "ssh: handshake failed: ssh: unable to authenticate, attempted methods [none password], no supported methods remain" {
//			t.Errorf("gossha: Wrong error: %s", err)
//		} else {
//			t.Logf("%s : error reported as expected", err)
//		}
//	} else {
//		t.Error("Error have to be thrown!")
//	}
//}
//
//func TestSendMessage(t *testing.T) {
//	t.Skipf("not working, need research")
//	//t.Parallel()
//	time.Sleep(200 * time.Millisecond)
//	client, session, err := connect("a", "a", testPort)
//	if err != nil {
//		t.Errorf("Connection error %s", err)
//		return
//	}
//	t.Logf("User a connected with password a")
//	err = session.Run("some test message\r")
//	if err != nil {
//		t.Errorf("Sending message error %s", err)
//		return
//	}
//	t.Logf("Command send")
//	err = session.Close()
//	if err != nil {
//		t.Errorf("%s : while  closing session", err)
//	}
//	t.Logf("Session closed")
//	err = client.Close()
//	if err != nil {
//		t.Errorf("%s : while closing clinet", err)
//	}
//	t.Logf("User a/a closed client")
//}
