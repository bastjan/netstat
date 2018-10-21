package netstat_test

import (
	"os/user"
	"testing"

	"gotest.tools/assert"

	"github.com/bastjan/netstat"
)

func TestConnectionUser(t *testing.T) {
	validUser, _ := user.Current()
	subject := &netstat.Connection{UserID: validUser.Uid}

	expectedUser, expectedErr := user.LookupId(validUser.Uid)
	user, err := subject.User()
	assert.DeepEqual(t, expectedUser, user)
	assert.DeepEqual(t, expectedErr, err)
}
