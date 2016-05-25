package config

import (
	"fmt"
	"os"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	//	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var sep = string(os.PathSeparator)

// Config is the object, that stories application parameters, populated from
// config file from `/etc/gossha/gossha.json`,`$HOME/.gossha/gossha.json` environment
// values prepended by prefix `GOSSHA_` or flags.
type Config struct {
	Port                    int
	Debug                   bool
	BindTo                  []string
	Driver                  string
	ConnectionString        string
	SSHPublicKeyPath        string
	SSHPrivateKeyPath       string
	Homedir                 string
	ExecuteOnMessage        string
	ExecuteOnPrivateMessage string
}

// Dump returns config as JSON object
func (c *Config) Dump() (string, error) {
	cfgTemplate := `# Automatically generated config file for GoSSHa - SSH powered chat
# Place it either in
#   /etc/gossha/gossha.toml
# or
#   ~/.gossha/gossha.toml
#
# Lines starting with # are comments
# See for details on toml syntax - https://github.com/toml-lang/toml

# Enable debug
#Debug=true
# Or disable it (as default behaviour)
#Debug=false

# On what port to listen for all interfaces (like for 0.0.0.0 address)
Port = %v
# Default value is 27015

# What additional addresses to bind to
# BindTo = ["127.0.0.1:27014","0.0.0.0:27016"]

#Setting database connection - various possible combinations are shown

#SQLite3 with database in local file
#Driver = "sqlite3"
#ConnectionString = "/var/lib/gossha/gossha.db"

#SQLite3 with database in memory
#Driver = "sqlite3"
#ConnectionString = ":memory:"

#MySQL database
#Driver = "mysql"
#ConnectionString = "username:password@hostname/database?charset=utf8&parseTime=True&loc=Local"

#PostgreSQL database. 1st variant
#Driver = "postgres"
#ConnectionString ="user=gorm dbname=gorm sslmode=disable"

#PostgreSQL database. 2nd variant
#Driver="postgres"
#ConnectionString="postgres://pqgotest:password@localhost/pqgotest?sslmode=verify-full")"

#This database connection setting are being used
Driver = "%s"
ConnectionString = "%s"

#Path to SSH Public key example
#SshPublicKeyPath = "/var/lib/gossha/id_rsa.pub"
#Path to be used
SshPublicKeyPath = "%s"


#Path to SSH Private key example
#SshPrivateKeyPath = "/var/lib/gossha/id_rsa"
SshPrivateKeyPath = "%s"

#Directory to search for custom scripts
#Homedir = "/var/lib/gossha/bin"
Homedir = "%s"

#Script to be executed on each message
#ExecuteOnMessage="/var/lib/gossha/onEachMessage.sh"
ExecuteOnMessage="%s"

#Script to be execute on each private message
#ExecuteOnPrivateMessage="/var/lib/gossha/onEachPrivateMessage.sh"
ExecuteOnPrivateMessage="%s"
	`

	return fmt.Sprintf(cfgTemplate,
		c.Port,
		c.Driver,
		c.ConnectionString,
		c.SSHPublicKeyPath,
		c.SSHPrivateKeyPath,
		c.Homedir,
		c.ExecuteOnMessage,
		c.ExecuteOnPrivateMessage,
	), nil
}

var RuntimeConfig *Config

// GetHomeDir returns the current working directory of application,
// usually the $HOME/.gossha
func GetHomeDir() (string, error) {
	hmdr, err := homedir.Dir()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%v%v.gossha", hmdr, sep), nil
}

// GetDatabasePath returns the current sqlite database path of application,
// usually the $HOME/.gossha/gossha.db
func GetDatabasePath() (string, error) {
	hmdr, err := homedir.Dir()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%v%v.gossha%vgossha.db", hmdr, sep, sep), nil
}

// GetPrivateKeyPath returns the current private ssh key location, usually
// the `$HOME/.ssh/id_rsa` one
func GetPrivateKeyPath() (string, error) {
	hmdr, err := homedir.Dir()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%v%v.ssh%vid_rsa", hmdr, sep, sep), nil
}

// GetPublicKeyPath returns the current public ssh key location, usually
// the `$HOME/.ssh/id_rsa.pub` one
func GetPublicKeyPath() (string, error) {
	hmdr, err := homedir.Dir()
	if err != nil {
		return "", err
	}
	sep := string(os.PathSeparator)
	return fmt.Sprintf("%v%v.ssh%vid_rsa.pub", hmdr, sep, sep), nil
}

// InitConfig creates configuration object, either by loading Toml file from this
// places - `/etc/gossha/gossha.toml` or `$HOME/.gossha/gossha.toml`, or by populating values
// from environment (with prefix of GOSSHA_). If the config file is not in place,
// the directory and file are created. Additional info about .toml file syntax
// can be seen here - https://github.com/toml-lang/toml
func InitConfig() (Config, error) {
	viper.SetConfigName("gossha")

	viper.SetEnvPrefix("gossha")

	viper.AddConfigPath("/etc/gossha")
	viper.AddConfigPath("$HOME/.gossha")
	viper.SetConfigType("toml")

	viper.AutomaticEnv()
	config := Config{}

	viper.SetDefault("port", 27015)
	viper.SetDefault("driver", "sqlite3")
	viper.SetDefault("debug", false)

	hmdr, err := GetHomeDir()
	if err != nil {
		return config, err
	}
	viper.SetDefault("homedir", hmdr)

	dbPath, err := GetDatabasePath()
	if err != nil {
		return config, err
	}
	viper.SetDefault("connectionString", dbPath)

	sshPublicKeyPath, err := GetPublicKeyPath()
	if err != nil {
		return config, err
	}
	viper.SetDefault("sshPublicKeyPath", sshPublicKeyPath)

	sshPrivateKeyPath, err := GetPrivateKeyPath()
	if err != nil {
		return config, err
	}
	viper.SetDefault("sshPrivateKeyPath", sshPrivateKeyPath)

	listOfAddressesToBind := []string{}
	viper.SetDefault("bindto", listOfAddressesToBind)

	err = viper.ReadInConfig()
	config.Port = viper.GetInt("port")
	config.Debug = viper.GetBool("debug")
	config.Driver = viper.GetString("driver")
	config.ConnectionString = viper.GetString("connectionString")
	config.SSHPublicKeyPath = viper.GetString("sshPublicKeyPath")
	config.SSHPrivateKeyPath = viper.GetString("sshPrivateKeyPath")
	config.Homedir = viper.GetString("homedir")
	config.ExecuteOnMessage = viper.GetString("executeOnMessage")
	config.ExecuteOnPrivateMessage = viper.GetString("executeOnPrivateMessage")
	config.BindTo = viper.GetStringSlice("bindto")
	RuntimeConfig = &config

	if err != nil {
		if err.Error() == "open : no such file or directory" {
			err := os.MkdirAll(hmdr, 0700)
			if err != nil {
				return config, err
			}

			configFileName := fmt.Sprintf("%v%vgossha.toml", hmdr, string(os.PathSeparator))
			fmt.Printf("Creating configuration file at %v...\n\n", configFileName)
			file, err := os.Create(configFileName)
			if err != nil {
				return config, err
			}

			configData, err := config.Dump()
			if err != nil {
				return config, err
			}

			defer file.Close()
			_, err = file.Write([]byte(configData))
			if err != nil {
				return config, err
			}
			return config, nil
		}
		return config, err
	}
	return config, nil
}
