package handler

/*
 * Templating - used to pretty print data if needed
 */

import (
	"fmt"
	"time"

	"github.com/vodolaz095/gossha/models"
	//	"golang.org/x/crypto/ssh"
	//	"golang.org/x/crypto/ssh/terminal"
)

// PrintPrompt makes promt for current user
func (h *Handler) PrintPrompt() string {
	return fmt.Sprintf("[%v@%v(%v) %v]{%v}:", h.CurrentUser.Name, h.Hostname, h.IP, "*", time.Now().Format("2006-1-2 15:04:05"))
}

// PrintMessage prints message in format of [username@hostname(192.168.1.2) *]{2006-1-2 15:04:05}:Hello!
func (h *Handler) PrintMessage(m *models.Message, u *models.User) string {
	var online string
	if u.IsOnline() {
		online = "*"
	} else {
		online = "x"
	}
	return fmt.Sprintf("[%v@%v(%v) %v]{%v}:%v\r\n", u.Name, m.Hostname, m.IP, online, m.CreatedAt.Format("2006-1-2 15:04:05"), m.Message)
}

// PrintNotification pretty prints the Notification recieved by Nerve into terminal given
func (h *Handler) PrintNotification(n *models.Notification) error {
	msg := h.PrintMessage(&n.Message, &n.User)
	//	if n.IsChat {
	//		_, err := h.writeToUser("%s%s%s", string(terminal.EscapeCodes.White), msg, string(terminal.EscapeCodes.Reset))
	//		return err
	//	}
	//	if n.IsSystem {
	//		_, err := h.writeToUser("%s%s%s", string(terminal.EscapeCodes.Green), msg, string(terminal.EscapeCodes.Reset))
	//		return err
	//	}
	//	if n.IsPrivateMessage {
	//		_, err := h.writeToUser("%s%s%s", string(terminal.EscapeCodes.Cyan), msg, string(terminal.EscapeCodes.Reset))
	//		return err
	//	}
	_, err := h.writeToUser(msg)
	return err
}
