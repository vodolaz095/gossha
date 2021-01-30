package handler

import (
	"fmt"
	"testing"
	"time"

	"github.com/vodolaz095/gossha/models"
)

func TestPrintPrompt(t *testing.T) {
	h := New()
	h.CurrentUser.Name = "test"
	h.IP = "8.8.8.8"
	prompt := h.PrintPrompt()
	wanted := fmt.Sprintf("{%s}[test@8.8.8.8 *]:", time.Now().Format(timeStampFormat))
	if prompt != wanted {
		t.Errorf("Wrong prompt:\nwanted: >>>%s<<<\ngot   :>>>%s<<<", wanted, prompt)
	}
}

func TestPrintMessage(t *testing.T) {
	h := New()
	h.CurrentUser.Name = "test"
	h.IP = "8.8.8.8"
	h.Hostname = "example.org"

	user := models.User{
		Name:           "test1",
		LastSeenOnline: time.Now(),
	}

	message := models.Message{
		Hostname:  "example.org",
		IP:        "8.8.8.8",
		CreatedAt: time.Now(),
		Message:   "hello",
	}
	prompt := h.PrintMessage(&message, &user)
	wanted := fmt.Sprintf("{%s}[test1@8.8.8.8 *]:hello", time.Now().Format(timeStampFormat))
	if prompt != wanted {
		t.Errorf("Wrong message:\nwanted: >>>%s<<<\ngot   :>>>%s<<<", wanted, prompt)
	}
}

type writerMock struct{}

var contents string

func (w writerMock) ReadPassword(prompt string) (line string, err error) {
	return "lalala", nil
}

func (w writerMock) SetPrompt(prompt string) {

}

func (w writerMock) Write(data []byte) (int, error) {
	contents = contents + string(data)
	return len(data), nil
}

func TestPrintNotification(t *testing.T) {
	term := writerMock{}
	h := New()
	h.Term = term
	h.CurrentUser.Name = "test"
	h.IP = "8.8.8.8"
	h.Hostname = "example.org"

	user := models.User{
		Name:           "test1",
		LastSeenOnline: time.Now(),
	}

	message := models.Message{
		Hostname:  "example.org",
		IP:        "8.8.8.8",
		CreatedAt: time.Now(),
		Message:   "hello",
	}

	notification := models.Notification{
		Message: message,
		User:    user,
	}

	err := h.PrintNotification(&notification)
	if err != nil {
		t.Errorf("Error writing to mock terminal %s", err)
	}
	prompt := contents
	test := fmt.Sprintf(
		"{%s}[test1@8.8.8.8 *]:hello\n\r", time.Now().Format(timeStampFormat))
	if prompt != test {
		t.Errorf("Wrong response of\n*%s*\ninstead of\n*%s*", prompt, test)
	}
}
