package netstat_test

import (
	"strings"
	"testing"

	"github.com/bastjan/netstat"
)

func TestTCPStateString(t *testing.T) {
	if netstat.TCPUnknown.String() != "" {
		t.Error("String() on UNKNOWN state should return empty string.")
	}
	if s := netstat.TCPListen.String(); s != "LISTEN" {
		t.Error("String() on TCPListen should return 'LISTEN'. Got:", s)
	}
	if s := netstat.TCPState(50).String(); !strings.Contains(s, "50") {
		t.Error("String() on invalid state should include it's number. Got: ", s)
	}
}
