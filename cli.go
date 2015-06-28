package gossha

import (
	"code.google.com/p/gopass"
	"encoding/base64"
	"fmt"
	"os"
	"strings"
)

// Greet writes some motivating text
func Greet() string {
	g := "ICBfX19fICAgICAgX19fXyBfX19fICBfICAgXyAgICAgICAKIC8gX19ffCBfX18vIF9fXy8gX19ffHwgfCB8IHwgX18gXyAKfCB8ICBfIC8gXyBcX19fIFxfX18gXHwgfF98IHwvIF9gIHwKfCB8X3wgfCAoXykgfF9fKSB8X18pIHwgIF8gIHwgKF98IHwKIFxfX19ffFxfX18vX19fXy9fX19fL3xffCB8X3xcX18sX3wKICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAK"
	gg, _ := base64.StdEncoding.DecodeString(g)
	ggg := make([]string, 0)
	ggg = append(ggg, string(gg))
	ggg = append(ggg, "Persistent SSH based chat for the ones, who cares.")
	ggg = append(ggg, fmt.Sprintf("Build: %v", VERSION))
	ggg = append(ggg, fmt.Sprintf("Version: %v", SUBVERSION))
	return strings.Join(ggg, "\r\n")
}

// PrintHelpOnCli prints help on console commands usage
func PrintHelpOnCli() {
	fmt.Println(Greet())
	fmt.Println("")
	fmt.Println("Console commands avaible: ")
	fmt.Println(" $ gossha ban [username] - delete user and all his/her messages")
	//	fmt.Println(" $ gossha log - list last 100 messages")
	fmt.Println(" $ gossha passwd [username] - create/update ordinary user by name and password")
	fmt.Println(" $ gossha root [username] - create/update root user by name and password")
	//	fmt.Println(" $ gossha users - print active users")
	fmt.Println("\nEmpty argument - start in server mode")
	fmt.Println("")
}

//ProcessConsoleCommand is a dispatcher for processing console commands,
//set by arguments used to start application
func ProcessConsoleCommand(a []string) {
	switch a[0] {
	case "-h":
		PrintHelpOnCli()
		os.Exit(1)
		break
	case "--help":
		PrintHelpOnCli()
		os.Exit(1)
		break
	case "root":
		if len(a) == 2 {
			name := a[1]
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
		break
	case "passwd":
		if len(a) == 2 {
			name := a[1]
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
		break
	case "users":
		fmt.Printf("Active users:\n")
		os.Exit(0)
		break
	case "ban":
		if len(a) == 2 {
			name := a[1]
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
		break
		//	case "log":
		//		fmt.Println("Recent messages:")
		//		h := New()
		//		n, err := h.GetMessages(100)
		//		if err != nil {
		//			panic(err)
		//		}
		//		for _, v := range n {
		//			fmt.Println(h.PrintNotification(&v))
		//		}
		//		fmt.Println("Thats all!")
		//		os.Exit(0)
		//		break
	default:
		os.Exit(1)
	}
}
