# v1.0.0 / 3.11.2018

Changes:
* 9affdbe Add go mod files
* 5df161a Add Connection.User() to lookup user associated with the connection
* 2ae5436 Add UserID, Transmit/ReceiveQueue to Connection
* 644fb01 Add TCPState TCPNewSynRecv
* a509188 Add Connection.State to represent the tcp state of a connection

Not Backwards Compatible Changes:
* 94a7073 Rename Netstat to Protocol, add Name field, reference in Connection
* 3dc8bec Connections(): Return pointers to connections
* a8c2beb Change type of Netstat from string to struct.

# v0.2.0-beta.1 / 16.10.2018

Not Backwards Compatible Changes:
* 73dd025 Rename Entries() to Connections() to be more symmetrical

# v0.1.0-beta.1 / 16.10.2018

Initial release
