package handler

import (
	"fmt"
	//	"reflect"
	"sort"
	"strings"

	//	"github.com/jinzhu/gorm"
	"github.com/vodolaz095/gossha/lib"
	"github.com/vodolaz095/gossha/models"
	"golang.org/x/crypto/ssh"
)

//TerminalInterface is interface repserenting terminal
type TerminalInterface interface {
	ReadPassword(prompt string) (line string, err error)
	Write([]byte) (numberOfBytesWritten int, err error)
}

//SSHChannelInterface is interface representing SSH channel session
type SSHChannelInterface interface {
	Close() error
}

// Board is a map of Handler's, with SessionId's as keys -
// http://godoc.org/golang.org/x/crypto/ssh#ConnMetadata
var Board map[string]*Handler

//KnownCommand includes command known to Handler and it's short description
type KnownCommand struct {
	Command     string
	Description string
}

// Handler is controller, that represents current user session
type Handler struct {
	SessionID          string
	KeyFingerPrint     ssh.PublicKey
	IP                 string
	Hostname           string
	LastShownMessageID int64
	CurrentUser        models.User
	Nerve              chan models.Notification
	KnownCommands      map[string]KnownCommand
	Term               TerminalInterface
	Connection         SSHChannelInterface
}

// New creates new Handler, representing users session
func New() Handler {
	n := make(chan models.Notification, 100)
	h := Handler{Nerve: n}
	h.KnownCommands = make(map[string]KnownCommand)
	h.addKnownCommand("h", "PrintHelpForUser", "(H)elp, show this screen")
	h.addKnownCommand("e", "Leave", "Close current session")
	h.addKnownCommand("q", "Leave", "Close current session")
	h.addKnownCommand("quit", "Leave", "Close current session")
	h.addKnownCommand("exit", "Leave", "Close current session")
	h.addKnownCommand("w", "Who", "List users, (W)ho are active on this server")
	h.addKnownCommand("i", "Info", "Print (I)nformation about yourself")
	h.addKnownCommand("k", "ImportPublicKey", "Use locally available SSH (K)eys to authorise your logins on this server")
	h.addKnownCommand("f", "ForgotPublicKey", "(F)orgot local available SSH key used for authorising your logins via this client")
	h.addKnownCommand("b", "Ban", "(B)an user (you need to have `root` permissions!)")
	h.addKnownCommand("r", "SignUpUser", "(R)egister new user (you need to have `root` permissions!)")
	h.addKnownCommand("rr", "SignUpRoot", "(R)egister new (r)oot user (you need to have `root` permissions!)")
	h.addKnownCommand("x", "ExecCommand", "E(X)ecutes custom user script from home directory")
	h.addKnownCommand("passwd", "ChangePassword", "Changes current user password")
	return h
}

func (h *Handler) writeToUser(format string, a ...interface{}) (bytesWriten int, err error) {
	return fmt.Fprintf(h.Term, format+"\n\r", a...)
}

func (h *Handler) addKnownCommand(key, commandName, help string) {
	h.KnownCommands[key] = KnownCommand{commandName, help}
}

// PrintHelpForUser outputs help for current user
func (h *Handler) PrintHelpForUser(command []string) error {

	h.writeToUser("GoSSHa - SSH powered chat. See https://github.com/vodolaz095/gossha for details...")
	//	h.writeToUser("Build #%v", VERSION)
	//	h.writeToUser("Version: %v", SUBVERSION)
	h.writeToUser("Commands available:")

	var keys []string
	for k, v := range h.KnownCommands {
		keys = append(keys, fmt.Sprintf(" \\%v - %v", k, v.Description))
	}
	sort.Strings(keys)
	for _, kk := range keys {
		h.writeToUser(kk)
	}
	h.writeToUser(" all other input is treated as message, that you send to server.")
	h.writeToUser("")
	h.writeToUser("")
	return nil
}

/*
 * Function to process user input in terminal
 * and call commands depending on input
 */

// AutoCompleteCallback is called for each keypress with
// the full input line and the current position of the cursor (in
// bytes, as an index into |line|). If it returns ok=false, the key
// press is processed normally. Otherwise it returns a replacement line
// and the new cursor position.
func (h *Handler) AutoCompleteCallback(s string, pos int, r rune) (string, int, bool) {
	//return s, pos, false
	//todo - under construction
	if string(r) == "\t" {
		//fmt.Printf("key:[%c] pos[%d] line:[%s]\n", r, pos, s)
		tokens := strings.Split(s, "")
		if len(tokens) > 0 {
			switch tokens[0] {
			case "\\":
				// todo - implement adding commands?
				//				fmt.Println("Looks like command!")
				break
			case "@":
				n := strings.TrimLeft(s, "@")
				if len(n) >= 1 {
					namePart := fmt.Sprint(n, "%")
					var user models.User
					err := models.DB.Table("user").Where("name LIKE ?", namePart).First(&user).Error
					//					fmt.Println("Looks like name!")
					if err == nil {
						return fmt.Sprintf("@%v", user.Name), (1 + len(user.Name)), true
					}
				}
				break
			default:
				//				fmt.Println("Looks like message!")
			}
		}
	}
	return s, pos, false
}

// ProcessCommand reads commands from terminal and processes them
func (h *Handler) ProcessCommand(command string) {
	tokens := strings.Split(command, "")
	if len(tokens) > 0 {
		switch tokens[0] {
		case "\\":
			a := strings.Split(strings.TrimLeft(command, "\\"), " ")
			f, ok := h.KnownCommands[a[0]]
			if ok {
				go func() { //todo not sure about it, need more tests
					err := lib.Invoke(h, f.Command, a)
					if err != nil {
						h.writeToUser("Error executing %v - %v", f.Command, err.Error())
					}
				}()
			} else {
				h.PrintHelpForUser(a)
			}
			break
		case "@":
			h.SendPrivateMessage(command)
			break
		default:
			err := h.SendMessage(command)
			if err != nil {
				h.writeToUser("Error sending message - %v", err.Error())
			}
		}
	}
}
