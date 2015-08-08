package gossha

import (
	"code.google.com/p/gopass"
	"encoding/base64"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

// Greet writes some motivating text alongside with application version
func Greet() string {
	g := "ICBfX19fICAgICAgX19fXyBfX19fICBfICAgXyAgICAgICAKIC8gX19ffCBfX18vIF9fXy8gX19ffHwgfCB8IHwgX18gXyAKfCB8ICBfIC8gXyBcX19fIFxfX18gXHwgfF98IHwvIF9gIHwKfCB8X3wgfCAoXykgfF9fKSB8X18pIHwgIF8gIHwgKF98IHwKIFxfX19ffFxfX18vX19fXy9fX19fL3xffCB8X3xcX18sX3wKICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAK"
	gg, _ := base64.StdEncoding.DecodeString(g)
	var ggg []string
	ggg = append(ggg, string(gg))
	ggg = append(ggg, "GoSSHa is a cross-platform ssh-server based chat program, with data persisted into relational databases of MySQL, PostgreSQL or Sqlite3. Public channel (with persisted messages) and private message (not stored) are supported. Application has serious custom scripting and hacking potential.")
	ggg = append(ggg, fmt.Sprintf("Build: %v", VERSION))
	ggg = append(ggg, fmt.Sprintf("Version: %v", SUBVERSION))
	ggg = append(ggg, "Homepages: https://github.com/vodolaz095/gossha")
	//	ggg = append(ggg, "           https://bitbucket.com/vodolaz095/gossha")
	ggg = append(ggg, "           https://godoc.com/github.com/vodolaz095/gossha")
	return strings.Join(ggg, "\r\n")
}

//ProcessConsoleCommand is a dispatcher for processing console commands,
//set by arguments used to start application
func ProcessConsoleCommand() {
	var rootCmd = &cobra.Command{
		Use: "gossha",
	}
	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Outputs program version and exits",
		Long:  "Outputs program version and exits",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
	var passwdCmd = &cobra.Command{
		Use:   "passwd [username]",
		Short: "Creates user or set new password to existent one",
		Long:  "Creates user or set new password to existent one",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 2 {
				name := args[1]
				password, err := gopass.GetPass("Enter password:")
				if err != nil {
					panic(err)
				}
				err = CreateUser(name, password, false)
				if err != nil {
					panic(err)
				}
				fmt.Printf("User %v is created and/or new password is set!\n", name)
				os.Exit(0)
			} else {
				fmt.Printf("Enter user's name!\n")
				os.Exit(1)
			}
		},
	}
	var makeRootUserCmd = &cobra.Command{
		Use:   "root [username]",
		Short: "Creates root user or set new password to existent one",
		Long:  "Creates root user or set new password to existent one",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 2 {
				name := args[1]
				password, err := gopass.GetPass("Enter password:")
				if err != nil {
					panic(err)
				}
				err = CreateUser(name, password, true)
				if err != nil {
					panic(err)
				}
				fmt.Printf("User %v is created and/or new password is set!\n", name)
				os.Exit(0)
			} else {
				fmt.Printf("Enter user's name!\n")
				os.Exit(1)
			}
		},
	}
	var banCmd = &cobra.Command{
		Use:   "ban [username]",
		Short: "Delete user and all his/her messages",
		Long:  "Delete user and all his/her messages",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 2 {
				name := args[1]
				err := BanUser(name)
				if err != nil {
					panic(err)
				}
				fmt.Printf("User %v is banned!\n", name)
				os.Exit(0)
			} else {
				fmt.Printf("Enter user's name!\n")
				os.Exit(1)
			}
		},
	}
	rootCmd.AddCommand(versionCmd, passwdCmd, makeRootUserCmd, banCmd)
	rootCmd.Execute()
}
