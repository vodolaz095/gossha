package lib

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
