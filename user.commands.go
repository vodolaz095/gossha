package gossha

/*
 * Commands related to users
 */

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"strings"
	"time"
)

// Leave notifies, that user has gone and close connection
// Handler is removed from Board in `ssh.go` file
func (h *Handler) Leave(connection ssh.Channel, term *terminal.Terminal, args []string) error {
	//delete(Board, h.SessionId)
	mesg := Message{
		//Id:        0,
		IP:        h.IP,
		Hostname:  h.Hostname,
		UserID:    h.CurrentUser.ID,
		Message:   "gone offline!",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	h.Broadcast(&mesg, true, false)
	term.Write([]byte("Goodbye!\n\r"))
	connection.Close()
	return nil
}

// Who lists the current active users
func (h *Handler) Who(connection ssh.Channel, term *terminal.Terminal, args []string) error {
	var cmds []string
	cmds = append(cmds, "Active sessions:\n\r")
	k := 0
	for _, v := range Board {
		k++
		cmds = append(cmds, fmt.Sprintf("%v) [%v@%v(%v) %v] {%v} ", k, v.CurrentUser.Name, v.Hostname, v.IP, v.CurrentUser.IsOnline(), v.CurrentUser.LastSeenOnline.Format("15:04:05")))
	}
	cmds = append(cmds, "\n\r\n\r")
	term.Write([]byte(strings.Join(cmds, "")))
	return nil
}

// Info prints additional information about yourself
func (h *Handler) Info(connection ssh.Channel, term *terminal.Terminal, args []string) error {
	var cmds []string
	var sessions []Session

	cmds = append(cmds, fmt.Sprintf("You are %v, logged from %v with IP of %v\n\r", h.CurrentUser.Name, h.Hostname, h.IP))
	cmds = append(cmds, "Your previous sessions: \n\r")

	err := DB.Table("session").Find(&sessions).Where("userId=?", h.CurrentUser.ID).Error
	if err != nil {
		return err
	}
	k := 0
	for _, v := range sessions {
		k++
		cmds = append(cmds, fmt.Sprintf("%v) %v(%v) since %v \n\r", k, v.Hostname, v.IP, v.CreatedAt.Format("15:04:05")))
	}
	cmds = append(cmds, " \n\r")
	term.Write([]byte(strings.Join(cmds, "")))
	return nil
}

// SignUpUser creates new user account, it requires root permissions
func (h *Handler) SignUpUser(connection ssh.Channel, term *terminal.Terminal, args []string) error {
	if h.CurrentUser.Root {
		fmt.Println(args)
		if len(args) == 2 {
			name := args[1]
			term.Write([]byte("Trying to signup/register/set password for user with name " + name + "!\n\rEnter new password:\n\r"))
			password, err := term.ReadPassword(">")
			if err != nil {
				return err
			}
			return CreateUser(name, password, false)
		}
		return fmt.Errorf("Name is empty, try `\\s someUserName`!")
	}
	return fmt.Errorf("You have to be root to signing up/registering/changing password!")
}

// SignUpRoot creates new user account, it requires root permissions
func (h *Handler) SignUpRoot(connection ssh.Channel, term *terminal.Terminal, args []string) error {
	if h.CurrentUser.Root {
		fmt.Println(args)
		if len(args) == 2 {
			name := args[1]
			term.Write([]byte("Trying to signup/register/set password for root user with name " + name + "!\n\rEnter new password:\n\r"))
			password, err := term.ReadPassword(">")
			if err != nil {
				return err
			}
			return CreateUser(name, password, true)
		}
		return fmt.Errorf("Name is empty, try `\\s someUserName`!")
	}
	return fmt.Errorf("You have to be root to signing up/registering/changing password!")
}

// Ban blocks user account, that is extracted from args, it requires root permissions
func (h *Handler) Ban(connection ssh.Channel, term *terminal.Terminal, args []string) error {
	if h.CurrentUser.Root {
		if len(args) == 2 {
			name := args[1]
			term.Write([]byte("Trying to ban " + name + "!\n\r"))
			return BanUser(name)
		}
		return fmt.Errorf("Name is empty, try `\\s someUserName`!")
	}
	return fmt.Errorf("You have to be root to signing up/registering/changing password!")
}

// ChangePassword sets the new password for current user
func (h *Handler) ChangePassword(connection ssh.Channel, term *terminal.Terminal, args []string) error {
	old, err := term.ReadPassword("Enter your old password:")
	if err != nil {
		return err
	}
	if h.CurrentUser.CheckPassword(old) {
		new1, err := term.ReadPassword("Enter your new password:")
		if err != nil {
			return err
		}
		new2, err := term.ReadPassword("Repeat your new password:")
		if err != nil {
			return err
		}
		if new1 == new2 {
			if len(new1) > 0 {
				term.Write([]byte("Setting new password...\r\n"))
				return h.CurrentUser.SetPassword(new1)
			}
			return fmt.Errorf("Unable to use empty password!")
		}
		return fmt.Errorf("Passwords do not match!")
	}
	return fmt.Errorf("Wrong password!")
}
