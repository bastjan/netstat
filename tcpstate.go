//go:generate stringer -type=TCPState -linecomment

package netstat

import (
	"strconv"
)

// TCPState represents the state of a TCP connection
// https://git.kernel.org/pub/scm/linux/kernel/git/torvalds/linux.git/tree/include/net/tcp_states.h?id=HEAD#n16.
type TCPState int

const (
	// TCPUnknown is an unknown state, 00 in /proc/net.
	TCPUnknown TCPState = iota //

	TCPEstablished // ESTABLISHED
	TCPSynSent     // SYN_SENT
	TCPSynRecv     // SYN_RECV
	TCPFinWait1    // FIN_WAIT1
	TCPFinWait2    // FIN_WAIT2
	TCPTimeWait    // TIME_WAIT
	TCPClose       // CLOSE
	TCPCloseWait   // CLOSE_WAIT
	TCPLastAck     // LAST_ACK
	TCPListen      // LISTEN
	TCPClosing     // CLOSING
	TCPNewSynRecv  // NEW_SYN_RECV
)

func tcpStatefromHex(hex string) TCPState {
	state, _ := strconv.ParseUint(hex, 16, 8)
	return TCPState(state)
}
