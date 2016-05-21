package ssh

import (
	"fmt"
	"net"

	"github.com/vodolaz095/gossha/handler"
	"github.com/vodolaz095/gossha/models"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

////HandlerInterface describes the handler
//type HandlerInterface interface {
//	MakeSSHConfig() (*ssh.ServerConfig, error)
//	PrintPrompt() string
//	AutoCompleteCallback(s string, pos int, r rune) (string, int, bool)
//	PrintHelpForUser(connection ssh.Channel, term *terminal.Terminal, command []string) error
//	GetMessages(int limit) ([]models.Notification, error)
//	ProcessCommand(connection ssh.Channel, term *terminal.Terminal, command []string) error
//	Leave(connection ssh.Channel, term *terminal.Terminal, command []string) error
//}

// StartSSHD starts the ssh server on address:port provided
func StartSSHD(addr string) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		panic(fmt.Sprintf("Failed to listen on %v port. Reason: (%s)", addr, err.Error()))
	}

	fmt.Printf("GoSSHa is listening on %v port!\n", addr)

	for {
		tcpConn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Failed to accept incoming connection (%s)\n", err)
			continue
		}
		h := handler.New()
		config := h.MakeSSHConfig()
		_, chans, reqs, err := ssh.NewServerConn(tcpConn, config)
		//		sshConn, chans, reqs, err := ssh.NewServerConn(tcpConn, config)
		if err != nil {
			fmt.Printf("Failed to handshake (%s)\n", err.Error())
			continue
		}

		//		fmt.Sprintf("New SSH connection from %s (%s)", sshConn.RemoteAddr(), sshConn.ClientVersion())
		go ssh.DiscardRequests(reqs)
		go handleChannels(chans, &h)
	}
}

func handleChannels(chans <-chan ssh.NewChannel, handler *Handler) {
	for newChannel := range chans {
		go handleChannel(newChannel, handler)
	}
}

func handleChannel(newChannel ssh.NewChannel, handler *Handler) {
	if t := newChannel.ChannelType(); t != "session" {
		newChannel.Reject(ssh.UnknownChannelType, fmt.Sprintf("unknown channel type: %s", t))
		return
	}
	connection, requests, err := newChannel.Accept()
	if err != nil {
		fmt.Printf("Could not accept channel (%s)", err)
		return
	}

	//http://play.golang.org/p/uN46-Pvd4O
	// Sessions have out-of-band requests such as "shell", "pty-req" and "env"
	go func() {
		for req := range requests {
			switch req.Type {
			case "pty-req":
				req.Reply(true, nil)
				break
			case "env":
				req.Reply(true, nil)
				break
			case "shell":
				// We only accept the default shell
				// (i.e. no command in the Payload)
				if len(req.Payload) == 0 {
					req.Reply(true, nil)
				}
				break
			default:
				req.Reply(true, nil)
			}
		}
	}()
	/*
		//http://play.golang.org/p/uN46-Pvd4O
		// Sessions have out-of-band requests such as "shell", "pty-req" and "env"
		//we try to utilize it
		go func() {
			for req := range requests {
				switch req.Type {
				case "shell":
					if len(req.Payload) == 0 {
						fmt.Println("Normal login.")
						req.Reply(true, nil)
					} else {
						cmd := fmt.Sprintf("Trying to execute command via `shell`: %v\n", string(req.Payload))
						connection.Write([]byte(cmd))
						fmt.Println(cmd)
						req.Reply(false, nil)
					}
					break
				case "exec":
					cmd := fmt.Sprintf("Trying to execute command via `exec`: %v\n", string(req.Payload))
					connection.Write([]byte(cmd))
					req.Reply(false, nil)
					fmt.Println(cmd)
					break
				default:
					req.Reply(true, nil)
				}
			}
		}()
	*/

	Board[handler.SessionID] = handler
	term := terminal.NewTerminal(connection, handler.PrintPrompt())
	term.AutoCompleteCallback = handler.AutoCompleteCallback
	handler.PrintHelpForUser(connection, term, []string{})
	msgs, err := handler.GetMessages(100)
	if err != nil {
		panic(err)
	}
	for _, v := range msgs {
		handler.PrintNotification(&v, term)
	}
	go func() {
		for {
			n1 := <-handler.Nerve
			handler.PrintNotification(&n1, term)
		}
	}()
	go func() {
		defer func() {
			connection.Close()
			delete(Board, handler.SessionID)
			handler.Leave(connection, term, []string{})
		}()
		for {
			line, err := term.ReadLine()
			if err != nil {
				break
			}
			handler.ProcessCommand(connection, term, line)
		}
	}()
}
