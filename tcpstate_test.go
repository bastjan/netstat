package netstat_test

import (
	"testing"

	"gotest.tools/assert"

	"github.com/bastjan/netstat"
)

func TestTCPStateString(t *testing.T) {
	assert.Equal(t, netstat.TCPUnknown.String(), "", "String() on UNKNOWN state should return empty string.")
	assert.Equal(t, netstat.TCPListen.String(), "LISTEN")
	assert.Equal(t, netstat.TCPState(50).String(), "TCPState(50)")
}
