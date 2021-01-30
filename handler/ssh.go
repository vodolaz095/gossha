package handler

/*
 * Authorization callbacks for SSH server by password and public key
 */

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/vodolaz095/gossha/config"
	"github.com/vodolaz095/gossha/models"
	"golang.org/x/crypto/ssh"
)

var authLog *log.Logger

// LoginByUsernameAndPassword is a authorization callback for ssh config
// see http://godoc.org/golang.org/x/crypto/ssh#ServerConfig for details
func (h *Handler) LoginByUsernameAndPassword(c ssh.ConnMetadata, password string) error {
	user := models.User{}
	name := c.User()
	ip := c.RemoteAddr().String()
	//hostname, err := lib.GetRemoteHostname(ip)
	//if err != nil {
	//	return err
	//}

	if err := models.DB.Table("user").Where("name=?", name).First(&user).Error; err == gorm.ErrRecordNotFound {
		return fmt.Errorf("user %v not found", name)
	}
	good, err := user.CheckPassword(password)
	if err != nil {
		return err
	}
	if good {
		h.SessionID = string(c.SessionID())
		h.CurrentUser = user
		h.IP = ip
		//h.Hostname = hostname
		h.CurrentUser.LastSeenOnline = time.Now()
		mesg := models.Message{
			IP:        h.IP,
			Hostname:  h.Hostname,
			UserID:    h.CurrentUser.ID,
			Message:   "appeared online!",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		h.Broadcast(&mesg, true, false)
		session := models.Session{
			UserID: user.ID,
			IP:     ip,
			//Hostname:  hostname,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		err = models.DB.Table("session").Save(&session).Error
		if err != nil {
			return err
		}
		return models.DB.Table("user").Save(&user).Error
	}
	return fmt.Errorf("wrong password for user %v", name)
}

// LoginByPublicKey is a authorization callback for ssh config
// see http://godoc.org/golang.org/x/crypto/ssh#ServerConfig for details
func (h *Handler) LoginByPublicKey(c ssh.ConnMetadata, publicKey string) error {
	key := models.Key{}
	user := models.User{}
	ip := c.RemoteAddr().String()
	//hostname, err := lib.GetRemoteHostname(ip)
	//if err != nil {
	//	return err
	//}
	err := models.DB.Table("key").Where("content=?", publicKey).First(&key).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("unknown public key")
		}
		return err
	}
	err = models.DB.Table("user").Where("id=? AND name = ?", key.UserID, c.User()).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("user %v not found", c.User())
		}
		return err
	}

	h.SessionID = string(c.SessionID())
	h.CurrentUser = user
	h.IP = ip
	//h.Hostname = hostname
	h.CurrentUser.LastSeenOnline = time.Now()

	mesg := models.Message{
		//Id:        0,
		IP:        h.IP,
		Hostname:  h.Hostname,
		UserID:    h.CurrentUser.ID,
		Message:   "appeared online!",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	h.Broadcast(&mesg, true, false)

	session := models.Session{
		UserID: user.ID,
		IP:     ip,
		//Hostname:  hostname,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err = models.DB.Table("session").Save(&session).Error
	if err != nil {
		return err
	}
	return models.DB.Table("user").Save(&user).Error
}

// MakeSSHConfig generates SSH server config used to authorize users
// to this handler context
// see http://godoc.org/golang.org/x/crypto/ssh#ServerConfig for details
func (h *Handler) MakeSSHConfig() *ssh.ServerConfig {
	authLog = log.New(os.Stdout, "[AUTH]", log.LstdFlags)

	sshConfig := &ssh.ServerConfig{
		PasswordCallback: func(c ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {
			return nil, h.LoginByUsernameAndPassword(c, string(pass))
		},
		PublicKeyCallback: func(c ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
			k := ssh.MarshalAuthorizedKey(key)
			h.KeyFingerPrint = key
			return nil, h.LoginByPublicKey(c, models.Hash(k))
		},
		AuthLogCallback: func(c ssh.ConnMetadata, method string, err error) {
			if err == nil {
				authLog.Printf("User %v@%v connected via %v\n", c.User(), c.RemoteAddr(), method)
			} else {
				authLog.Printf("User %v@%v failed to connect via %v, because %v\n", c.User(), c.RemoteAddr(), method, err.Error())
			}
		},
		NoClientAuth: false,
	}
	privateBytes, err := ioutil.ReadFile(config.RuntimeConfig.SSHPrivateKeyPath)
	if err != nil {
		panic("Failed to load private key from " + config.RuntimeConfig.SSHPrivateKeyPath)
	}
	private, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		panic("Failed to parse private key")
	}

	sshConfig.AddHostKey(private)
	return sshConfig
}
