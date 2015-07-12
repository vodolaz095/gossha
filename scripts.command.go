package gossha

import (
	"bufio"
	"fmt"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"io/ioutil"
	"os"
	"os/exec"
	"time"
)

// ExecCommand executes custom user scripts from home directory. It checks, if
//the script file is executable, and than it executes it, setting the environment
// parameters with this values:
func (h *Handler) ExecCommand(connection ssh.Channel, term *terminal.Terminal, input []string) error {
	if len(input) == 2 {
		cmd := input[1]
		mesg := Message{
			IP:        h.IP,
			Hostname:  h.Hostname,
			UserID:    h.CurrentUser.ID,
			Message:   fmt.Sprintf("Tries to execute command '%v'", cmd),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		err := DB.Table("Message").Create(&mesg).Error
		if err != nil {
			return err
		}
		h.Broadcast(&mesg, true, false)
		scriptPath := fmt.Sprintf("%v%vscripts%v%v", GetHomeDir(), sep, sep, cmd)
		fi, err := os.Stat(scriptPath)
		if err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("Command %v is not known!", cmd)
			}
			return err
		}

		if fi.IsDir() {
			return fmt.Errorf("Unable to execute command, the %v is a directory!", scriptPath)
		}
		if fi.Mode().Perm()&0111 != 0 {
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
			if h.CurrentUser.Root {
				err = os.Setenv("GOSSHA_ROOT", "true")
			} else {
				err = os.Setenv("GOSSHA_ROOT", "false")
			}
			if err != nil {
				return err
			}
			term.Write([]byte(fmt.Sprintf("Executing command '%v':\r\n\r\n", cmd)))
			chld := exec.Command(scriptPath)
			out, err := chld.StdoutPipe()
			chld.Start()
			defer func() {
				s := chld.Wait() // Doesn't block
				if s != nil {
					term.Write([]byte(fmt.Sprintf("\r\nCommand '%v' failed with status code %v!!!\r\n", cmd, s)))
				} else {
					term.Write([]byte(fmt.Sprintf("\r\nCommand '%v' executed properly!!!\r\n", cmd)))
				}

			}()
			if err != nil {
				return err
			}
			scn := bufio.NewScanner(out)

			for scn.Scan() {
				term.Write([]byte(fmt.Sprintf("%v\r\n", scn.Text())))
			}
			if err := scn.Err(); err != nil {
				return err
			}
			return nil
		}
		return fmt.Errorf("Unable to execute command, the %v is not a executable!", scriptPath)
	}
	files, err := ioutil.ReadDir(fmt.Sprintf("%v%vscripts", GetHomeDir(), sep))
	if err != nil {
		return err
	}
	term.Write([]byte(fmt.Sprintf("Commands avaible:\r\n")))
	for _, v := range files {
		if v.Mode().Perm()&0111 != 0 {
			term.Write([]byte(fmt.Sprintf("  %v\r\n", v.Name())))
		}
	}
	return nil
}
