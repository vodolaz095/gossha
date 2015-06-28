package gossha

/*
 * Authorization callbacks for SSH server by password and public key
 */

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"strings"
	"time"
)

// LoginByUsernameAndPassword is a authorization callback for ssh config
// see http://godoc.org/golang.org/x/crypto/ssh#ServerConfig for details
func (h *Handler) LoginByUsernameAndPassword(c ssh.ConnMetadata, password string) error {
	user := User{}
	name := c.User()
	ip := strings.Split(c.RemoteAddr().String(), ":")[0]
	hostname, err := GetRemoteHostname(ip)
	if err != nil {
		return err
	}

	if err := DB.Table("user").Where("name=?", name).First(&user).Error; err == gorm.RecordNotFound {
		return fmt.Errorf("User %v not found!", name)
	}

	if user.CheckPassword(password) {
		h.SessionId = string(c.SessionID())
		h.CurrentUser = user
		h.Ip = ip
		h.Hostname = hostname
		h.CurrentUser.LastSeenOnline = time.Now()
		mesg := Message{
			Ip:        h.Ip,
			Hostname:  h.Hostname,
			UserID:    h.CurrentUser.ID,
			Message:   "appeared online!",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		h.Broadcast(&mesg, true, false)
		session := Session{
			UserID:    user.ID,
			Ip:        ip,
			Hostname:  hostname,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		err = DB.Table("session").Save(&session).Error
		if err != nil {
			return err
		}
		return DB.Table("user").Save(&user).Error
	} else {
		return fmt.Errorf("Wrong password for user %v!", name)
	}
}

// LoginByPublicKey is a authorization callback for ssh config
// see http://godoc.org/golang.org/x/crypto/ssh#ServerConfig for details
func (h *Handler) LoginByPublicKey(c ssh.ConnMetadata, publicKey string) error {
	key := Key{}
	user := User{}
	ip := strings.Split(c.RemoteAddr().String(), ":")[0]
	hostname, err := GetRemoteHostname(ip)
	if err != nil {
		return err
	}
	err = DB.Table("key").Where("content=?", publicKey).First(&key).Error
	if err != nil {
		if err == gorm.RecordNotFound {
			return fmt.Errorf("Public key is not known!")
		} else {
			return err
		}
	}
	err = DB.Table("user").Where("id=? AND name = ?", key.UserID, c.User()).First(&user).Error
	if err != nil {
		if err == gorm.RecordNotFound {
			return fmt.Errorf("User %v not found!", c.User())
		} else {
			return err
		}
	}

	h.SessionId = string(c.SessionID())
	h.CurrentUser = user
	h.Ip = ip
	h.Hostname = hostname
	h.CurrentUser.LastSeenOnline = time.Now()

	mesg := Message{
		//Id:        0,
		Ip:        h.Ip,
		Hostname:  h.Hostname,
		UserID:    h.CurrentUser.ID,
		Message:   "appeared online!",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	h.Broadcast(&mesg, true, false)

	session := Session{
		UserID:    user.ID,
		Ip:        ip,
		Hostname:  hostname,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err = DB.Table("session").Save(&session).Error
	if err != nil {
		return err
	}
	return DB.Table("user").Save(&user).Error
}

// MakeSSHConfig generates SSH server config used to authorize users
// to this handler context
// see http://godoc.org/golang.org/x/crypto/ssh#ServerConfig for details
func (h *Handler) MakeSSHConfig() *ssh.ServerConfig {
	config := &ssh.ServerConfig{
		PasswordCallback: func(c ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {
			return nil, h.LoginByUsernameAndPassword(c, string(pass))
		},
		PublicKeyCallback: func(c ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
			k := string(ssh.MarshalAuthorizedKey(key))
			//fmt.Printf("Public key is %v -- %v\n", key.Type(), k)
			h.KeyFingerPrint = key
			return nil, h.LoginByPublicKey(c, Hash(k))
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
	privateBytes, err := ioutil.ReadFile(RuntimeConfig.SshPrivateKeyPath)
	if err != nil {
		panic("Failed to load private key from " + RuntimeConfig.SshPrivateKeyPath)
	}
	private, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		panic("Failed to parse private key")
	}

	config.AddHostKey(private)
	return config
}
