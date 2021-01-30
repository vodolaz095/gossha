package handler

/*
 * Templating - used to pretty print data if needed
 */

import (
	"fmt"
	"github.com/vodolaz095/gossha/models"
	"time"
	//	"golang.org/x/crypto/ssh"
	//"golang.org/x/crypto/ssh/terminal"
)

// PrintPrompt makes promt for current user
func (h *Handler) PrintPrompt() string {
	return fmt.Sprintf(
		"{%v}[%s@%s *]:",
		time.Now().Format(timeStampFormat), h.CurrentUser.Name, h.IP,
	)
}

// PrintMessage prints message in format of [username@192.168.1.2:33921 *]{2006-1-2 15:04:05}:Hello!
func (h *Handler) PrintMessage(m *models.Message, u *models.User) string {
	var online string
	if u.IsOnline() {
		online = "*"
	} else {
		online = "x"
	}
	return fmt.Sprintf(
		"{%s}[%s@%s %s]:%s",
		m.CreatedAt.Format(timeStampFormat),
		u.Name, m.IP, online, m.Message,
	)
}

// PrintNotification pretty prints the Notification received by Nerve into terminal given
func (h *Handler) PrintNotification(n *models.Notification) error {
	msg := h.PrintMessage(&n.Message, &n.User)
	_, err := h.writeToUser(msg)
	return err
}
