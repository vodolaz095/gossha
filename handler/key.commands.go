package handler

/*
 * User commands to process ssh keys fingerprints
 */

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/vodolaz095/gossha/models"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

// ImportPublicKey saves hash of public key from Handler.KeyFingerPrint
// into `key` table
func (h *Handler) ImportPublicKey(connection ssh.Channel, term *terminal.Terminal, input []string) error {
	fingerprint := h.KeyFingerPrint
	if fingerprint == nil {
		return fmt.Errorf("Public key is empty!")
	}
	term.Write([]byte("Importing public key...\n\r"))
	k := string(ssh.MarshalAuthorizedKey(fingerprint))
	key := models.Key{
		UserID:  h.CurrentUser.ID,
		Content: models.Hash(k),
	}
	err := models.DB.Table("key").Create(&key).Error
	if err != nil {
		term.Write([]byte("Error importing key, probably it is already imported!\n\r"))
		return err
	}
	term.Write([]byte("Key imported succesefully!\n\r"))
	return nil
}

// ForgotPublicKey removes corresponding public key fingerprint hash from `key` table
func (h *Handler) ForgotPublicKey(connection ssh.Channel, term *terminal.Terminal, input []string) error {
	fingerprint := h.KeyFingerPrint
	if fingerprint == nil {
		return fmt.Errorf("Public key is empty!")
	}
	k := string(ssh.MarshalAuthorizedKey(fingerprint))
	key := models.Key{
		UserID:  h.CurrentUser.ID,
		Content: models.Hash(k),
	}
	err := models.DB.Table("key").Where("content=?", models.Hash(k)).First(&key).Error
	if err == nil {
		err = models.DB.Table("key").Delete(&key).Error
		if err != nil {
			return err
		}
		term.Write([]byte("Public key is removed! You'll need password in future to be able to authorize from this client.\n\r"))
		return nil
	}
	if err == gorm.ErrRecordNotFound {
		return nil
	}
	return err
}
