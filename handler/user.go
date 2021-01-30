package handler

/*
 * Commands related to users
 */

import (
	"fmt"
	"time"

	"github.com/vodolaz095/gossha/models"
)

// Leave notifies, that user has gone and close connection
// Handler is removed from Board in `ssh.go` file
func (h *Handler) Leave(args []string) error {
	mesg := models.Message{
		//Id:        0,
		IP:        h.IP,
		Hostname:  h.Hostname,
		UserID:    h.CurrentUser.ID,
		Message:   "gone offline!",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	h.Broadcast(&mesg, true, false)
	h.writeToUser("Goodbye!")
	return h.Connection.Close()
}

// Who lists the current active users
func (h *Handler) Who(args []string) error {
	h.writeToUser("Active sessions:")
	k := 0
	for _, v := range Board {
		k++
		if v.CurrentUser.Name != "" {
			h.writeToUser("%d) [%s@%s [%s] %t] {%s} ",
				k, v.CurrentUser.Name, v.Hostname, v.IP, v.CurrentUser.IsOnline(), v.CurrentUser.LastSeenOnline.Format("15:04:05"),
			)
		}

	}
	h.writeToUser("")
	h.writeToUser("")
	return nil
}

// Info prints additional information about yourself
func (h *Handler) Info(args []string) error {
	var sessions []models.Session

	h.writeToUser("You are %v, logged from %v with IP of %v.", h.CurrentUser.Name, h.Hostname, h.IP)
	h.writeToUser("Your previous sessions: ")

	err := models.DB.Table("session").Find(&sessions).Where("userId=?", h.CurrentUser.ID).Error
	if err != nil {
		return err
	}
	k := 0
	for _, v := range sessions {
		k++
		h.writeToUser("%v) %v(%v) since %v \n\r", k, v.Hostname, v.IP, v.CreatedAt.Format("15:04:05"))
	}
	h.writeToUser("")
	return nil
}

// SignUpUser creates new user account, it requires root permissions
func (h *Handler) SignUpUser(args []string) error {
	if h.CurrentUser.Root {
		//fmt.Println(args)
		switch len(args) {
		case 3:
			name := args[1]
			password := args[2]
			err := models.CreateUser(name, password, false)
			if err != nil {
				return err
			}
			_, err = h.writeToUser("Password for %s is set!", name)
			return err
			//TODO - not working, because `h.Term.ReadPassword` behaves unconsistent
			//		case 2:
			//			name := args[1]
			//			//			fmt.Println("Waiting for password")
			//			//			h.Term.Write([]byte(fmt.Sprintf("\nEnter password for ordinary user %s:\n", name)))
			//			password, err := h.Term.ReadPassword(fmt.Sprintf("Enter password for user `%s`:", name))
			//			h.Term.SetPrompt(h.PrintPrompt())
			//			if err != nil {
			//				return err
			//			}
			//			//			fmt.Println("Password recieved", password)
			//			//			fmt.Println("Prompt updated, creating user")
			//			return models.CreateUser(name, password, false)

		default:
			return fmt.Errorf("Try `\\r someUserName [newPassword]` to sign up or change password for somebody!")
		}
	}
	return fmt.Errorf("You have to be root to signing up/registering/changing password!")
}

// SignUpRoot creates new user account, it requires root permissions
func (h *Handler) SignUpRoot(args []string) error {
	if h.CurrentUser.Root {
		//fmt.Println(args)
		switch len(args) {
		case 3:
			name := args[1]
			password := args[2]
			err := models.CreateUser(name, password, true)
			if err != nil {
				return err
			}
			_, err = h.writeToUser("Password for %s is set!", name)
			return err

			//TODO - not working, because `h.Term.ReadPassword` behaves unconsistent
			//		case 2:
			//			name := args[1]
			//			h.Term.Write([]byte(fmt.Sprintf("\nEnter password for ADMIN user %s:\n", name)))
			//			password, err := h.Term.ReadPassword(fmt.Sprintf("Enter password for root `%s`:", name))
			//			h.Term.SetPrompt(h.PrintPrompt())
			//			if err != nil {
			//				return err
			//			}
			//			return models.CreateUser(name, password, true)

		default:
			return fmt.Errorf("Try `\\rr someUserName newPassword` to sign up or change password for somebody with root permissions!")
		}
	}
	return fmt.Errorf("You have to be root to signing up/registering/changing password!")
}

// Ban blocks user account, that is extracted from args, it requires root permissions
func (h *Handler) Ban(args []string) error {
	if h.CurrentUser.Root {
		if len(args) == 2 {
			name := args[1]
			h.writeToUser("Trying to ban %s!", name)
			return models.BanUser(name)
		}
		return fmt.Errorf("Name is empty, try `\\b someUserName`!")
	}
	return fmt.Errorf("You have to be root to signing up/registering/changing password!")
}

// ChangePassword sets the new password for current user
func (h *Handler) ChangePassword(args []string) error {
	old, err := h.Term.ReadPassword("Enter your old password:")
	if err != nil {
		return err
	}
	good, err := h.CurrentUser.CheckPassword(old)
	if err != nil {
		return err
	}
	if good {
		new1, err := h.Term.ReadPassword("Enter your new password:")
		if err != nil {
			return err
		}
		new2, err := h.Term.ReadPassword("Repeat your new password:")
		if err != nil {
			return err
		}
		if new1 == new2 {
			if len(new1) > 0 {
				_, err1 := h.writeToUser("Setting new password...")
				if err1 != nil {
					return err1
				}
				return h.CurrentUser.SetPassword(new1)
			}
			return fmt.Errorf("unable to use empty password!")
		}
		return fmt.Errorf("passwords do not match")
	}
	return fmt.Errorf("wrong password")
}
