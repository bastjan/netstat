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

		UserID:     "6523",
		IP:         net.ParseIP("127.0.0.1"),
		Port:       38911,
		RemoteIP:   net.ParseIP("0.0.0.0"),
		RemotePort: 0,
		State:      netstat.TCPListen,

		TransmitQueue: 50,
		ReceiveQueue:  100,
	}

	tcp6Connection = netstat.Connection{
		Exe:     "",
		Cmdline: []string{"/usr/bin/bundle", "exec", "puma", "-p41703"},
		Pid:     3002,

		Inode: 44365,

		UserID:     "6523",
		IP:         net.ParseIP("2001::4:0:131b"),
		Port:       41703,
		RemoteIP:   net.ParseIP("::"),
		RemotePort: 0,
		State:      netstat.TCPListen,
	}
)

func init() {
	netstat.ProcRoot = "./test/proc"
}

func TestConnections(t *testing.T) {
	connections, err := netstat.TCP.Connections()
	assert.NilError(t, err)
	assert.DeepEqual(t, connections, []netstat.Connection{tcpConnection})

	connections, err = netstat.TCP6.Connections()
	assert.NilError(t, err)
	assert.DeepEqual(t, connections, []netstat.Connection{tcp6Connection})
}

func TestConnectionsProcNetNotFound(t *testing.T) {
	_, err := (&netstat.Netstat{RelPath: "./nothere"}).Connections()
	assert.ErrorContains(t, err, "can't open proc file")
	assert.ErrorContains(t, err, "test/proc/nothere")
}

func TestConnectionsEmptyFileDoesNotCrashNetstat(t *testing.T) {
	_, err := (&netstat.Netstat{RelPath: "net/empty"}).Connections()
	assert.ErrorContains(t, err, "net/empty has no content")
}
