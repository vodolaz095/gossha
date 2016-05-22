package handler

/*
 * Authorization callbacks for SSH server by password and public key
 */

import (
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/vodolaz095/gossha/config"
	"github.com/vodolaz095/gossha/lib"
	"github.com/vodolaz095/gossha/models"
	"golang.org/x/crypto/ssh"
)

// LoginByUsernameAndPassword is a authorization callback for ssh config
// see http://godoc.org/golang.org/x/crypto/ssh#ServerConfig for details
func (h *Handler) LoginByUsernameAndPassword(c ssh.ConnMetadata, password string) error {
	user := models.User{}
	name := c.User()
	ip := strings.Split(c.RemoteAddr().String(), ":")[0]
	hostname, err := lib.GetRemoteHostname(ip)
	if err != nil {
		return err
	}

	if err := models.DB.Table("user").Where("name=?", name).First(&user).Error; err == gorm.ErrRecordNotFound {
		return fmt.Errorf("User %v not found!", name)
	}

	if user.CheckPassword(password) {
		h.SessionID = string(c.SessionID())
		h.CurrentUser = user
		h.IP = ip
		h.Hostname = hostname
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
			UserID:    user.ID,
			IP:        ip,
			Hostname:  hostname,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		err = models.DB.Table("session").Save(&session).Error
		if err != nil {
			return err
		}
		return models.DB.Table("user").Save(&user).Error
	}
	return fmt.Errorf("Wrong password for user %v!", name)
}

// LoginByPublicKey is a authorization callback for ssh config
// see http://godoc.org/golang.org/x/crypto/ssh#ServerConfig for details
func (h *Handler) LoginByPublicKey(c ssh.ConnMetadata, publicKey string) error {
	key := models.Key{}
	user := models.User{}
	ip := strings.Split(c.RemoteAddr().String(), ":")[0]
	hostname, err := lib.GetRemoteHostname(ip)
	if err != nil {
		return err
	}
	err = models.DB.Table("key").Where("content=?", publicKey).First(&key).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("Public key is not known!")
		}
		return err
	}
	err = models.DB.Table("user").Where("id=? AND name = ?", key.UserID, c.User()).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("User %v not found!", c.User())
		}
		return err
	}

	h.SessionID = string(c.SessionID())
	h.CurrentUser = user
	h.IP = ip
	h.Hostname = hostname
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
		UserID:    user.ID,
		IP:        ip,
		Hostname:  hostname,
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
	sshConfig := &ssh.ServerConfig{
		PasswordCallback: func(c ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {
			return nil, h.LoginByUsernameAndPassword(c, string(pass))
		},
		PublicKeyCallback: func(c ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
			k := string(ssh.MarshalAuthorizedKey(key))
			//fmt.Printf("Public key is %v -- %v\n", key.Type(), k)
			h.KeyFingerPrint = key
			return nil, h.LoginByPublicKey(c, models.Hash(k))
		},
		AuthLogCallback: func(c ssh.ConnMetadata, method string, err error) {
			if err == nil {
				fmt.Printf("Connection success from %v@%v via %v\n", c.User(), c.RemoteAddr(), method)
			} else {
				fmt.Printf("Connection fail from %v@%v via %v. Reason: %v\n", c.User(), c.RemoteAddr(), method, err.Error())
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
