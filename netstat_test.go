package netstat_test

import (
	"net"
	"testing"

	"gotest.tools/assert"

	"github.com/bastjan/netstat"
)

var (
	tcpConnection = netstat.Connection{
		Exe:     "/bin/sleep",
		Cmdline: []string{},
		Pid:     3001,

		Inode: 44360,

		IP:         net.ParseIP("127.0.0.1"),
		Port:       38911,
		RemoteIP:   net.ParseIP("0.0.0.0"),
		RemotePort: 0,
	}

	tcp6Connection = netstat.Connection{
		Exe:     "",
		Cmdline: []string{"/usr/bin/bundle", "exec", "puma", "-p41703"},
		Pid:     3002,

		Inode: 44365,

		IP:         net.ParseIP("2001::4:0:131b"),
		Port:       41703,
		RemoteIP:   net.ParseIP("::"),
		RemotePort: 0,
	}
)

func init() {
	netstat.ProcRoot = "./test/proc"
}

func TestEntries(t *testing.T) {
	entries, err := netstat.TCP.Entries()
	assert.NilError(t, err)
	assert.DeepEqual(t, entries, []netstat.Connection{tcpConnection})

	entries, err = netstat.TCP6.Entries()
	assert.NilError(t, err)
	assert.DeepEqual(t, entries, []netstat.Connection{tcp6Connection})
}

func TestEntriesProcNetNotFound(t *testing.T) {
	_, err := netstat.Netstat("./nothere").Entries()
	assert.ErrorContains(t, err, "can't open proc file")
	assert.ErrorContains(t, err, "test/proc/nothere")
}

func TestEntriesEmptyFileDoesNotCrashNetstat(t *testing.T) {
	_, err := netstat.Netstat("net/empty").Entries()
	assert.ErrorContains(t, err, "net/empty has no content")
}
