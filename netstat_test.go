package netstat_test

import (
	"net"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bastjan/netstat"
)

var (
	tcpConnection = &netstat.Connection{
		Protocol: netstat.TCP,

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

	tcp6Connection = &netstat.Connection{
		Protocol: netstat.TCP6,

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
	compareResult := func(p *netstat.Protocol, expected []*netstat.Connection) {
		connections, err := p.Connections()
		if err != nil {
			t.Error("Connections() returned unexpected errors:", err)
		}
		if diff := cmp.Diff(connections, expected); diff != "" {
			t.Error("Connections() returned connections differ from expected connections:\n", diff)
		}
	}
	compareResult(netstat.TCP, []*netstat.Connection{tcpConnection})
	compareResult(netstat.TCP6, []*netstat.Connection{tcp6Connection})
}

func TestConnectionsProcNetNotFound(t *testing.T) {
	_, err := (&netstat.Protocol{RelPath: "./nothere"}).Connections()
	expectError(t, err, "test/proc/nothere: no such file or directory", "Connections() should return an error if the proc file can't be found")
}

func TestConnectionsEmptyFile(t *testing.T) {
	_, err := (&netstat.Protocol{RelPath: "net/empty"}).Connections()
	expectError(t, err, "net/empty has no content", "Connections() should return an error if net file is empty")
}

func expectError(t *testing.T, err error, expectedErr, nilMessage string) {
	t.Helper()
	if err == nil {
		t.Fatal(nilMessage)
	}
	if strings.Contains(err.Error(), expectedErr) {
		return
	}
	t.Error("Error message should contain filename and error.", "Expected:", expectedErr, "Got:", err.Error())
}
