package gossha

import (
	"fmt"
	//	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
	//	"reflect"
	"sort"
	"strings"
)

//KnownCommand includes command known to Handler and it's short description
type KnownCommand struct {
	Command string
	//	Command     func(ssh.Channel, *terminal.Terminal, ...string) error
	Description string
}

// Handler is controller, that represents current user session
type Handler struct {
	SessionID          string
	KeyFingerPrint     ssh.PublicKey
	IP                 string
	Hostname           string
	LastShownMessageID int64
	CurrentUser        User
	Nerve              chan Notification
	KnownCommands      map[string]KnownCommand
}

// New creates new Handler, representing users session
func New() Handler {
	n := make(chan Notification, 100)
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
	h.addKnownCommand("f", "ForgotPublicKey", "(F)orgot localy available SSH key used for authorising your logins via this client")
	h.addKnownCommand("b", "Ban", "(B)an user (you need to have `root` permissions!)")
	h.addKnownCommand("r", "SignUpUser", "(R)egisters new user (you need to have `root` permissions!)")
	h.addKnownCommand("rr", "SignUpRoot", "(R)egisters new (r)oot user (you need to have `root` permissions!)")
	h.addKnownCommand("x", "ExecCommand", "E(X)ecutes custom user script from home directory")
	h.addKnownCommand("passwd", "ChangePassword", "Changes current user password")
	return h
}

func (h *Handler) addKnownCommand(key, commandName, help string) {
	h.KnownCommands[key] = KnownCommand{commandName, help}
}

// PrintHelpForUser outputs help for current user
func (h *Handler) PrintHelpForUser(connection ssh.Channel, term *terminal.Terminal, command []string) error {
	var cmds []string
	cmds = append(cmds, "GoSSHa - very secure chat.\n\r")
	cmds = append(cmds, fmt.Sprintf("Build #%v \n\r", VERSION))
	cmds = append(cmds, fmt.Sprintf("Version: %v \n\r", SUBVERSION))
	cmds = append(cmds, fmt.Sprintf("Commands avaible:\n\r"))
	var keys []string
	for k, v := range h.KnownCommands {
		keys = append(keys, fmt.Sprintf(" \\%v - %v\n\r", k, v.Description))
	}
	sort.Strings(keys)
	for _, kk := range keys {
		cmds = append(cmds, kk)
	}
	cmds = append(cmds, " all other input is treated as message, that you send to server\n\r")
	cmds = append(cmds, " \n\r")
	cmds = append(cmds, " \n\r")
	term.Write([]byte(strings.Join(cmds, "")))
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
					var user User
					err := DB.Table("user").Where("name LIKE ?", namePart).First(&user).Error
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
func (h *Handler) ProcessCommand(connection ssh.Channel, term *terminal.Terminal, command string) {
	tokens := strings.Split(command, "")
	switch tokens[0] {
	case "\\":
		a := strings.Split(strings.TrimLeft(command, "\\"), " ")
		f, ok := h.KnownCommands[a[0]]
		if ok {
			go func() { //todo not sure about it, need more tests
				err := Invoke(h, f.Command, connection, term, a)
				if err != nil {
					term.Write([]byte(fmt.Sprintf("Error - %v\n\r", err.Error())))
				}
			}()
		} else {
			h.PrintHelpForUser(connection, term, a)
		}
		break
	case "@":
		h.SendPrivateMessage(connection, term, command)
		break
	default:
		err := h.SendMessage(connection, term, command)
		if err != nil {
			term.Write([]byte(fmt.Sprintf("Error - %v\n\r", err.Error())))
		}

	}
}
