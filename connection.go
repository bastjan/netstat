package netstat

import (
	"net"
	"os/user"
)

// Connection contains the gathered information about an open network connection.
type Connection struct {
	// Exe contains the path to the process.
	// Exe is empty if there was an error reading /proc/pid/exe.
	Exe string
	// Cmdline contains the complete command line for the process split by \000. Trailing \000 removed.
	// Returns an empty array if /proc/pid/cmdline can't be read.
	Cmdline []string
	// Pid contains the pid of the process. Is zero if open connection can't be assigned to a pid.
	Pid int

	// UserID represents the user account id of the user owning the socket.
	// On Linux systems it is usually a uint32.
	// Type string was chosen because os/user.LookupId() wants a string.
	UserID string

	// Inode contains the inode for the open connection.
	Inode uint64

	// IP holds the local IP for the connection.
	IP net.IP
	// Port holds the local port for the connection.
	Port int
	// RemoteIP holds the remote IP for the connection.
	RemoteIP net.IP
	// RemotePort holds the remote port for the connection.
	RemotePort int
	// State represents the state of a TCP connection. The UDP 'states' shown
	// are recycled from TCP connection states and have a slightly different meaning.
	State TCPState

	// TransmitQueue is the outgoing data queue in terms of kernel memory usage in bytes.
	TransmitQueue uint64
	// ReceiveQueue is the incoming data queue in terms of kernel memory usage in bytes.
	ReceiveQueue uint64

	// Protocol contains the protocol this connection was discovered with.
	Protocol *Protocol
}

// User looks up the user owning the socket.
// If the user cannot be found, the returned error is of type UnknownUserIdError.
func (c *Connection) User() (*user.User, error) {
	return user.LookupId(c.UserID)
}
