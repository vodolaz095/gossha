package ssh

import (
	"fmt"
	"github.com/vodolaz095/gossha/handler"
	"log"
	"net"
	"os"
	//	"github.com/vodolaz095/gossha/models"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

// Good read - http://play.golang.org/p/uN46-Pvd4O

var sshdLog *log.Logger

// StartSSHD starts the ssh server on address:port provided
func StartSSHD(addr string) error {
	sshdLog = log.New(os.Stdout, "[SSHD]", log.LstdFlags)
	handler.Board = make(map[string]*handler.Handler, 0)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("%s - while binding to listen on %v port", err, addr)
	}
	sshdLog.Printf("GoSSHa is listening on %v port!\n", addr)
	for {
		tcpConn, err := listener.Accept()
		if err != nil {
			sshdLog.Printf("Failed to accept incoming connection (%s)\n", err)
			continue
		}
		h := handler.New()
		config := h.MakeSSHConfig()
		_, chans, reqs, err := ssh.NewServerConn(tcpConn, config)
		//		sshConn, chans, reqs, err := ssh.NewServerConn(tcpConn, config)
		if err != nil {
			sshdLog.Printf("Failed to handshake (%s)\n", err.Error())
			continue
		}

		//		fmt.Sprintf("New SSH connection from %s (%s)", sshConn.RemoteAddr(), sshConn.ClientVersion())
		go ssh.DiscardRequests(reqs)
		go handleChannels(chans, &h)
	}
}

func handleChannels(chans <-chan ssh.NewChannel, h *handler.Handler) {
	for newChannel := range chans {
		go handleChannel(newChannel, h)
	}
}

func handleChannel(newChannel ssh.NewChannel, h *handler.Handler) {
	if t := newChannel.ChannelType(); t != "session" {
		newChannel.Reject(ssh.UnknownChannelType, fmt.Sprintf("unknown channel type: %s", t))
		return
	}
	connection, requests, err := newChannel.Accept()
	if err != nil {
		fmt.Printf("Could not accept channel (%s)", err)
		return
	}
	go func() {
		for req := range requests {
			switch req.Type {
			case "pty-req":
				req.Reply(true, nil)
				break
			case "env":
				req.Reply(true, nil)
				break
			case "exec":
				cmd := string(req.Payload)
				sshdLog.Printf("Trying to execute command %s via `exec`...", cmd)
				h.ProcessCommand(cmd)
				req.Reply(false, nil)
				break
			case "shell":
				cmd := string(req.Payload)
				if len(cmd) == 0 {
					req.Reply(true, nil)
				} else {
					sshdLog.Printf("Trying to execute command %s via `shell`...", cmd)
					h.ProcessCommand(cmd)
					req.Reply(false, nil)
				}
				break
			default:
				req.Reply(true, nil)
			}
		}
	}()
	handler.Board[h.SessionID] = h
	term := terminal.NewTerminal(connection, h.PrintPrompt())
	term.AutoCompleteCallback = h.AutoCompleteCallback
	h.Term = term
	h.Connection = connection
	h.PrintHelpForUser([]string{})
	msgs, err := h.GetMessages(100)
	if err != nil {
		panic(err)
	}
	for _, v := range msgs {
		h.PrintNotification(&v)
	}
	go func() {
		for {
			n1 := <-h.Nerve
			h.PrintNotification(&n1)
		}
	}()
	go func() {
		defer func() {
			connection.Close()
			delete(handler.Board, h.SessionID)
			h.Leave([]string{})
		}()
		for {
			line, err := term.ReadLine()
			if err != nil {
				break
			}
			h.ProcessCommand(line)
		}
	}()
}
