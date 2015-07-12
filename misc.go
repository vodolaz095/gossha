package gossha

/*
 * Assorted functions, that do not deserve to be a member of any struct
 */

import (
	"fmt"
	"net"
	"reflect"
	"strings"
)

// GetRemoteHostname tries to resolve remote hostname by ip
func GetRemoteHostname(a string) (string, error) {
	hostnames, err := net.LookupAddr(a)
	if err != nil {
		return "", nil
	}
	if hostnames[0] != "" {
		return hostnames[0], nil
	}
	return "localhost", nil
}

// MakeDSNHelp returns some help regarding database connection string
func MakeDSNHelp() string {
	var dsnHelpArr []string
	dsnHelpArr = append(dsnHelpArr, "Database connection string. Examples:")
	dsnHelpArr = append(dsnHelpArr, "   	--driver=sqlite3 --connectionString=/var/lib/gossha/gossha.db")
	dsnHelpArr = append(dsnHelpArr, "   	--driver=sqlite3 --connectionString=:memory:")
	dsnHelpArr = append(dsnHelpArr, "   	--driver=mysql --connectionString=user:password@localhost/dbname?charset=utf8&parseTime=True&loc=Local")
	dsnHelpArr = append(dsnHelpArr, "   	--driver=postgres --connectionString='user=gorm dbname=gorm sslmode=disable'")
	dsnHelpArr = append(dsnHelpArr, "   	--driver=postgres --connectionString=postgres://pqgotest:password@localhost/pqgotest?sslmode=verify-full")
	dsnHelp := strings.Join(dsnHelpArr, "\n")
	return dsnHelp
}

// Invoke calls struct method by name. Used in calling commands.
// See http://stackoverflow.com/questions/8103617/call-a-struct-and-its-method-by-name-in-go
func Invoke(any interface{}, name string, args ...interface{}) error {
	inputs := make([]reflect.Value, len(args))
	for i := range args {
		inputs[i] = reflect.ValueOf(args[i])
	}
	a := reflect.ValueOf(any).MethodByName(name).Call(inputs)
	err := a[0].Interface()
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	return nil
}
