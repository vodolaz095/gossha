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
func (h *Handler) ImportPublicKey(input []string) error {
	fingerprint := h.KeyFingerPrint
	if fingerprint == nil {
		return fmt.Errorf("empty public key")
	}
	_, err := h.writeToUser("Importing public key...")
	if err != nil {
		return err
	}
	k := ssh.MarshalAuthorizedKey(fingerprint)
	key := models.Key{
		UserID:  h.CurrentUser.ID,
		Content: models.Hash(k),
	}
	err = models.DB.
		Table("key").
		Create(&key).
		Error
	if err != nil {
		_, err = h.writeToUser("Error importing key, probably it is already imported!")
		return err
	}
	_, err = h.writeToUser("Key imported successfully!")
	handlerLog.Printf("User %s@%s imported public key %s",
		h.CurrentUser.Name,
		h.IP,
		key.Content,
	)

	return err
}

// ForgotPublicKey removes corresponding public key fingerprint hash from `key` table
func (h *Handler) ForgotPublicKey(connection ssh.Channel, term *terminal.Terminal, input []string) error {
	fingerprint := h.KeyFingerPrint
	if fingerprint == nil {
		return fmt.Errorf("public key is empty")
	}
	k := ssh.MarshalAuthorizedKey(fingerprint)
	key := models.Key{
		UserID:  h.CurrentUser.ID,
		Content: models.Hash(k),
	}
	err := models.
		DB.
		Table("key").
		Where("content=?", key.Content).
		First(&key).
		Error

	if err == nil {
		err = models.DB.Table("key").Delete(&key).Error
		if err != nil {
			return err
		}
		h.writeToUser("Public key is removed! You'll need password in future to be able to authorize from this client.")
		handlerLog.Printf("User %s@%s removed public key %s",
			h.CurrentUser.Name,
			h.IP,
			key.Content,
		)
		return nil
	}
	if err == gorm.ErrRecordNotFound {
		return nil
	}
	return err
}
