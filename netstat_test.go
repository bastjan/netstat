package netstat_test

import (
	"testing"

	"gotest.tools/assert"

	"github.com/bastjan/netstat"
)

var (
	TCP  = netstat.Netstat("./test/proc/net/tcp")
	TCP6 = netstat.Netstat("./test/proc/net/tcp6")
)

var (
	tcpEntry = netstat.Entry{
		Inode: 44360,

		IP:         "127.0.0.1",
		Port:       38911,
		RemoteIP:   "127.0.0.1",
		RemotePort: 0,
	}

	tcpEntry6 = netstat.Entry{
		Inode: 44365,

		IP:         "00:00:00:00:00:00:00:00",
		Port:       41703,
		RemoteIP:   "00:00:00:00:00:00:00:00",
		RemotePort: 0,
	}
)

func TestEntries(t *testing.T) {
	entries, err := TCP.Entries()
	assert.NilError(t, err)
	assert.DeepEqual(t, entries, []netstat.Entry{tcpEntry})

	entries, err = TCP6.Entries()
	assert.NilError(t, err)
	assert.DeepEqual(t, entries, []netstat.Entry{tcpEntry6})
}
