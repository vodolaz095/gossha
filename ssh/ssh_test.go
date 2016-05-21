package ssh

import (
	"testing"
)

const addressToListen = "0.0.0.0:33972"

func TestSSHServerWorks(t *testing.T) {
	t.Parallel()
	err := StartSSHD(addressToListen)
	if err != nil {
		t.Errorf("Error starting ssh server %s", err)
	}

}
