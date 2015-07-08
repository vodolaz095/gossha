package gossha

import (
	"encoding/json"
	"fmt"
	homedir "github.com/mitchellh/go-homedir"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"
)

var sep = string(os.PathSeparator)

// Config is the object, that stories application parameters, populated from
// config file from `/etc/gossha/gossha.json`,`$HOME/.gossha/gossha.json` environment
// values prepended by prefix `GOSSHA_` or flags.
type Config struct {
	Port                    int      `json:"port"`
	Debug                   bool     `json:"debug"`
	BindTo                  []string `json:"bindTo"`
	Driver                  string   `json:"driver"`
	ConnectionString        string   `json:"connectionString"`
	SSHPublicKeyPath        string   `json:"sshPublicKeyPath"`
	SSHPrivateKeyPath       string   `json:"sshPrivateKeyPath"`
	Homedir                 string   `json:"homedir"`
	ExecuteOnMessage        string   `json:"executeOnMessage"`
	ExecuteOnPrivateMessage string   `json:"executeOnPrivateMessage"`
}

// GetHomeDir returns the current working directory of application,
// usually the $HOME/.gossha
func GetHomeDir() string {
	hmdr, err := homedir.Dir()
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%v%v.gossha", hmdr, sep)
}

// GetDatabasePath returns the current sqlite database path of application,
// usually the $HOME/.gossha/gossha.db
func GetDatabasePath() string {
	hmdr, err := homedir.Dir()
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%v%v.gossha%vgossha.db", hmdr, sep, sep)
}

// GetPrivateKeyPath returns the current private ssh key location, usually
// the `$HOME/.ssh/id_rsa` one
func GetPrivateKeyPath() string {
	hmdr, err := homedir.Dir()
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%v%v.ssh%vid_rsa", hmdr, sep, sep)
}

// GetPublicKeyPath returns the current public ssh key location, usually
// the `$HOME/.ssh/id_rsa.pub` one
func GetPublicKeyPath() string {
	hmdr, err := homedir.Dir()
	if err != nil {
		panic(err)
	}
	sep := string(os.PathSeparator)
	return fmt.Sprintf("%v%v.ssh%vid_rsa.pub", hmdr, sep, sep)
}

// InitConfig creates configuration object, either by loading JSON file from this
// places - `/etc/gossha/gossha.json` or `$HOME/.gossha/gossha.json`, or by populating values
// from environment (with prefix of GOSSHA_). If the config file is not in place,
//the directory and file are created
func InitConfig() (Config, []string, error) {
	viper.SetConfigName("gossha")
	viper.AddConfigPath("/etc/gossha")
	viper.AddConfigPath("$HOME/.gossha")
	viper.SetConfigType("json") //todo - maybe not needed
	viper.SetEnvPrefix("gossha")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	config := Config{}

	flag.Uint("port", 27015, "set the port to listen for connections")
	viper.BindPFlag("port", flag.Lookup("port"))
	flag.Bool("debug", false, "start pprof debugging on port 3000")
	viper.BindPFlag("port", flag.Lookup("port"))

	flag.String("driver", "sqlite3", "set the database driver to use, possible values are `sqlite3`,`mysql`,`postgres`")
	viper.BindPFlag("driver", flag.Lookup("driver"))

	flag.String("connectionString", GetDatabasePath(), MakeDSNHelp())
	viper.BindPFlag("connectionString", flag.Lookup("connectionString"))

	flag.String("sshPublicKeyPath", GetPublicKeyPath(), "location of public ssh key to be used with server, usually the $HOME/.ssh/id_rsa.pub")
	viper.BindPFlag("sshPublicKeyPath", flag.Lookup("sshPublicKeyPath"))

	flag.String("sshPrivateKeyPath", GetPrivateKeyPath(), "location of private ssh key to be used with server, usually the $HOME/.ssh/id_rsa")
	viper.BindPFlag("sshPrivateKeyPath", flag.Lookup("sshPrivateKeyPath"))

	flag.String("homedir", GetHomeDir(), "The home directory of module, usually $HOME/.gossha")
	viper.BindPFlag("homedir", flag.Lookup("homedir"))

	flag.String("executeOnMessage", "", "Script to execute on each message")
	viper.BindPFlag("executeOnMessage", flag.Lookup("executeOnMessage"))

	flag.String("executeOnPrivateMessage", "", "Script to execute on each private message")
	viper.BindPFlag("executeOnPrivateMessage", flag.Lookup("executeOnPrivateMessage"))

	flag.Parse()
	config.Port = viper.GetInt("port")
	config.Debug = viper.GetBool("debug")
	config.Driver = viper.GetString("driver")
	config.ConnectionString = viper.GetString("connectionString")
	config.SSHPublicKeyPath = viper.GetString("sshPublicKeyPath")
	config.SSHPrivateKeyPath = viper.GetString("sshPrivateKeyPath")
	config.Homedir = viper.GetString("homedir")
	config.ExecuteOnMessage = viper.GetString("executeOnMessage")
	config.ExecuteOnPrivateMessage = viper.GetString("executeOnPrivateMessage")
	listOfAddressesToBind := []string{}
	viper.SetDefault("bindto", listOfAddressesToBind)
	config.BindTo = viper.GetStringSlice("bindto")

	if config.Debug {
		for k, v := range viper.AllSettings() {
			fmt.Printf("DEBUG: setting config %v=%v\n", k, v)
		}
	}

	if err != nil {
		if err.Error() == "open : no such file or directory" {

			hmdr := GetHomeDir()
			err := os.MkdirAll(hmdr, 0700)
			if err != nil {
				return config, flag.Args(), err
			}

			configFileName := fmt.Sprintf("%v%vgossha.json", config.Homedir, string(os.PathSeparator))
			fmt.Printf("Creating configuration file at %v...\n\n", configFileName)
			file, err := os.Create(configFileName)
			if err != nil {
				return config, flag.Args(), err
			}

			configData, err := config.dump()
			if err != nil {
				return config, flag.Args(), err
			}

			defer file.Close()
			_, err = file.Write([]byte(configData))
			if err != nil {
				return config, flag.Args(), err
			}
			return config, flag.Args(), nil
		}
		return config, flag.Args(), err

	}
	return config, flag.Args(), nil
}

func (c *Config) dump() (string, error) {
	data, err := json.MarshalIndent(c, "", "  ")
	return string(data), err
}
