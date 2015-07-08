package gossha

/*
 * User commands to process messages
 */

import (
	//"fmt"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"os/exec"
	"strings"
	"time"
)

// SendMessage sends message from this user into chat. Message is saved into persistent datastorage.
// Also the command from `RuntimeConfig.ExecuteOnMessage` is executed if present
func (h *Handler) SendMessage(connection ssh.Channel, term *terminal.Terminal, input string) error {
	//authorized
	var comment string
	if len(input) > 255 {
		comment = string([]byte(input)[0:255])
	} else {
		comment = input
	}

	mesg := Message{
		IP:        h.IP,
		Hostname:  h.Hostname,
		UserID:    h.CurrentUser.ID,
		Message:   comment,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := DB.Table("message").Create(&mesg).Error
	h.LastShownMessageID = mesg.ID
	if err != nil {
		return err
	}

	h.CurrentUser.LastSeenOnline = time.Now()
	h.Broadcast(&mesg, false, true)
	err = DB.Table("user").Save(&h.CurrentUser).Error
	if err != nil {
		return err
	}

	if RuntimeConfig.ExecuteOnMessage != "" {
		err = os.Setenv("GOSSHA_USERNAME", h.CurrentUser.Name)
		if err != nil {
			return err
		}
		err = os.Setenv("GOSSHA_HOSTNAME", h.Hostname)
		if err != nil {
			return err
		}
		err = os.Setenv("GOSSHA_IP", h.IP)
		if err != nil {
			return err
		}
		err = os.Setenv("GOSSHA_MESSAGE", comment)
		if err != nil {
			return err
		}

		if h.CurrentUser.Root {
			err = os.Setenv("GOSSHA_ROOT", "true")
			if err != nil {
				return err
			}
		} else {
			err = os.Setenv("GOSSHA_ROOT", "false")
			if err != nil {
				return err
			}
		}
		chld := exec.Command(RuntimeConfig.ExecuteOnMessage)
		_, err := chld.StdoutPipe()
		chld.Start()
		//fmt.Println("Executing", output)
		return err
	}
	return nil
}

//SendPrivateMessage delivers message to the reciever only, message is not saved into persistent datastorage
// Also the command from `RuntimeConfig.ExecuteOnPrivateMessage` is executed if present
func (h *Handler) SendPrivateMessage(connection ssh.Channel, term *terminal.Terminal, input string) error {
	if string(input[0]) != "@" {
		return nil
	}
	var to string
	messageSend := false
	tokens := strings.Split(strings.TrimLeft(input, "@"), " ")
	name := tokens[0]
	for _, v := range Board {
		if v.CurrentUser.Name == name {
			to = v.CurrentUser.Name
			mesg := Message{
				IP:        h.IP,
				Hostname:  h.Hostname,
				UserID:    h.CurrentUser.ID,
				Message:   input,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			//reciever can  see private messages he sends
			v.Nerve <- Notification{
				User:             h.CurrentUser,
				Message:          mesg,
				IsSystem:         false,
				IsChat:           false,
				IsPrivateMessage: true,
			}
			//also author can see private messages he/she sends
			h.Nerve <- Notification{
				User:             h.CurrentUser,
				Message:          mesg,
				IsSystem:         false,
				IsChat:           false,
				IsPrivateMessage: true,
			}
			messageSend = true
		}
	}
	if !messageSend {
		term.Write([]byte("Unable to send private message, user is offline!"))
	}
	if RuntimeConfig.ExecuteOnPrivateMessage != "" {
		err := os.Setenv("GOSSHA_USERNAME", h.CurrentUser.Name)
		if err != nil {
			return err
		}
		err = os.Setenv("GOSSHA_HOSTNAME", h.Hostname)
		if err != nil {
			return err
		}
		err = os.Setenv("GOSSHA_IP", h.IP)
		if err != nil {
			return err
		}
		err = os.Setenv("GOSSHA_MESSAGE", input)
		if err != nil {
			return err
		}

		if h.CurrentUser.Root {
			err = os.Setenv("GOSSHA_ROOT", "true")
			if err != nil {
				return err
			}
		} else {
			err = os.Setenv("GOSSHA_ROOT", "false")
			if err != nil {
				return err
			}
		}

		err = os.Setenv("GOSSHA_TO", to)
		if err != nil {
			return err
		}

		chld := exec.Command(RuntimeConfig.ExecuteOnPrivateMessage)
		_, err = chld.StdoutPipe()
		chld.Start()
		//fmt.Println("Executing", output)
		return err
	}
	return nil
}

// Broadcast sends Message in form of Notification
// to all other Handler's, each of thems corresponding authorized User.
func (h *Handler) Broadcast(mesg *Message, isSystem, isChat bool) {
	for k, v := range Board {
		if k != h.SessionID {
			v.Nerve <- Notification{
				User:             h.CurrentUser,
				Message:          *mesg,
				IsSystem:         isSystem,
				IsChat:           isChat,
				IsPrivateMessage: false,
			}
		}
	}
}

// PrivateMessage sends Message in form of Notification to all handlers, which have the
// Handler.CurrentUser.Name equal to first argument
func (h *Handler) PrivateMessage(name string, mesg *Message) {
	for k, v := range Board {
		if k != h.SessionID {
			if h.CurrentUser.Name == name {
				v.Nerve <- Notification{User: h.CurrentUser, Message: *mesg, IsSystem: false, IsChat: false, IsPrivateMessage: true}
			}
		}
	}
}

// GetMessages outputs recent messages in form of Notification array
func (h *Handler) GetMessages(limit int) ([]Notification, error) {
	ret := make([]Notification, 0)
	var messages []Message
	var l int64
	DB.Table("message").Preload("User").Where("message.id > ?", h.LastShownMessageID).Limit(limit).Order("message.id asc").Find(&messages)
	for _, m := range messages {
		ret = append(ret, Notification{User: m.User, Message: m, IsSystem: false, IsChat: true})
	}

	if len(messages) > 0 {
		h.LastShownMessageID = l
	}
	return ret, nil
}
