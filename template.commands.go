package gossha

/*
 * Templating - used to pretty print data if needed
 */

import (
	"fmt"
	"time"
	//	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

// PrintPrompt makes promt for current user
func (h *Handler) PrintPrompt() string {
	return fmt.Sprintf("[%v@%v(%v) %v]{%v}:", h.CurrentUser.Name, h.Hostname, h.IP, "*", time.Now().Format("2006-1-2 15:04:05"))
}

// PrintMessage prints message in format of [username@hostname(192.168.1.2) *]{2006-1-2 15:04:05}:Hello!
func (h *Handler) PrintMessage(m *Message, u *User) string {
	var online string
	if u.IsOnline() {
		online = "*"
	} else {
		online = "x"
	}
	return fmt.Sprintf("[%v@%v(%v) %v]{%v}:%v\r\n", u.Name, m.Hostname, m.IP, online, m.CreatedAt.Format("2006-1-2 15:04:05"), m.Message)
}

// PrintNotification pretty prints the Notification recieved by Nerve into terminal given
func (h *Handler) PrintNotification(n *Notification, term *terminal.Terminal) error {
	msg := []byte(h.PrintMessage(&n.Message, &n.User))
	if n.IsChat {
		_, err := term.Write(term.Escape.White)
		_, err = term.Write(msg)
		_, err = term.Write(term.Escape.Reset)
		return err
	}
	if n.IsSystem {
		_, err := term.Write(term.Escape.Green)
		_, err = term.Write(msg)
		_, err = term.Write(term.Escape.Reset)
		return err
	}
	if n.IsPrivateMessage {
		_, err := term.Write(term.Escape.Cyan)
		_, err = term.Write(msg)
		_, err = term.Write(term.Escape.Reset)
		return err
	}
	_, err := term.Write(msg)
	return err
}
