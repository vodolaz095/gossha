package lib

import (
	"fmt"
	"testing"
)

func TestGetRemoteHostname(t *testing.T) {
	ips, err := GetRemoteHostname("127.0.0.1")
	if err != nil {
		t.Errorf("Error getting hostname %s", err)
	}
	if ips != "localhost.localdomain." {
		t.Errorf("Wrong localhost %s instead of localhost.localdomain.", ips)
	}
}

type invokeTarget struct{}

func (i *invokeTarget) Do(a, b string) error {
	return fmt.Errorf("ok %s %s", a, b)
}

func TestInvoke(t *testing.T) {
	it := invokeTarget{}
	err := it.Do("da", "du")
	if err.Error() != "ok da du" {
		t.Errorf("Unable to call `Invoke` properly!")
	}
	//TODO - fix
	//	err = Invoke(it, "Do", "da", "du")

	//	if err.Error() != "ok da du" {
	//		t.Errorf("Unable to call `Invoke` properly!")
	//	}
}
