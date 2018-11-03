package netstat_test

import (
	"os/user"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bastjan/netstat"
)

func TestConnectionUser(t *testing.T) {
	validUser, _ := user.Current()
	subject := &netstat.Connection{UserID: validUser.Uid}

	expectedUser, expectedErr := user.LookupId(validUser.Uid)
	user, err := subject.User()
	if diff := cmp.Diff(expectedUser, user); diff != "" {
		t.Error("Connection.User() should return the same user as user.LookupId().", diff)
	}
	if diff := cmp.Diff(expectedErr, err); diff != "" {
		t.Error("Connection.User() should return the same error as user.LookupId().", diff)
	}
}
