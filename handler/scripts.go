package handler

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"time"

	"github.com/vodolaz095/gossha/config"
	"github.com/vodolaz095/gossha/models"
)

var sep = string(os.PathSeparator)

// ExecCommand executes custom user scripts from home directory. It checks, if
//the script file is executable, and than it executes it, setting the environment
// parameters with this values:
func (h *Handler) ExecCommand(input []string) error {
	if len(input) == 2 {
		cmd := input[1]
		mesg := models.Message{
			IP:        h.IP,
			Hostname:  h.Hostname,
			UserID:    h.CurrentUser.ID,
			Message:   fmt.Sprintf("Tries to execute command '%v'", cmd),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		err := models.DB.Table("Message").Create(&mesg).Error
		if err != nil {
			return err
		}
		h.Broadcast(&mesg, true, false)
		hmdr, err := config.GetHomeDir()
		if err != nil {
			return err
		}
		scriptPath := fmt.Sprintf("%v%vscripts%v%v", hmdr, sep, sep, cmd)
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
			h.writeToUser("Executing command '%v':", cmd)
			chld := exec.Command(scriptPath)
			out, err := chld.StdoutPipe()
			chld.Start()
			defer func() {
				s := chld.Wait() // Doesn't block
				if s != nil {
					h.writeToUser("Command '%v' failed with status code %v!!!", cmd, s)
				} else {
					h.writeToUser("Command '%v' executed properly!!!", cmd)
				}

			}()
			if err != nil {
				return err
			}
			scn := bufio.NewScanner(out)

			for scn.Scan() {
				h.writeToUser("%v", scn.Text())
			}
			if err := scn.Err(); err != nil {
				return err
			}
			return nil
		}
		return fmt.Errorf("Unable to execute command, the %v is not a executable!", scriptPath)
	}
	hmdr, err := config.GetHomeDir()
	if err != nil {
		return err
	}

	files, err := ioutil.ReadDir(fmt.Sprintf("%v%vscripts", hmdr, sep))
	if err != nil {
		return err
	}
	h.writeToUser("Commands available:")
	for _, v := range files {
		if v.Mode().Perm()&0111 != 0 {
			h.writeToUser("  %v", v.Name())
		}
	}
	return nil
}
