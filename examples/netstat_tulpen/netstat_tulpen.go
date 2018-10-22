/*
Sample program to the netstat package. Prints almost the same output as `netstat -tulpen`.

Differences I've seen:
- Some hints are missing e.g. Not all processes could be identified,...
- udp shows State CLOSE instead of nothing
*/
package main

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/bastjan/netstat"
)

var header = []string{"Proto", "Recv-Q", "Send-Q", "Local Address", "Foreign Address", "State", "User", "Inode", "PID/Program name"}

func main() {
	out := [][]string{header}
	out = append(out, get("tcp", netstat.TCP)...)
	out = append(out, get("tcp6", netstat.TCP6)...)
	out = append(out, get("udp", netstat.UDP)...)
	out = append(out, get("udp6", netstat.UDP6)...)

	printAligned(out)
}

func get(protoName string, loc *netstat.Netstat) [][]string {
	connections, _ := loc.Connections()
	results := make([][]string, 0, len(connections))
	for _, conn := range connections {
		if !isListening(protoName, conn.State) {
			continue
		}
		results = append(results, []string{
			protoName,
			strconv.FormatUint(conn.ReceiveQueue, 10),
			strconv.FormatUint(conn.TransmitQueue, 10),
			fmt.Sprintf("%s:%s", conn.IP, formatPort(conn.Port)),
			fmt.Sprintf("%s:%s", conn.RemoteIP, formatPort(conn.RemotePort)),
			conn.State.String(),
			conn.UserID,
			strconv.FormatUint(conn.Inode, 10),
			formatPidProgname(conn.Pid, conn.Exe),
		})
	}

	return results
}

func isListening(protoName string, tcpState netstat.TCPState) bool {
	tcpListen := strings.HasPrefix(protoName, "tcp") && tcpState == netstat.TCPListen
	udpListen := strings.HasPrefix(protoName, "udp") && tcpState == netstat.TCPClose
	return tcpListen || udpListen
}

func formatPort(port int) string {
	if port == 0 {
		return "*"
	}
	return strconv.Itoa(port)
}

func formatPidProgname(pid int, exe string) string {
	if pid == 0 {
		return "-"
	}
	_, binary := filepath.Split(exe)
	return fmt.Sprintf("%d/%s", pid, binary)
}

func printAligned(table [][]string) {
	widths := make([]int, len(table[0]))

	for _, row := range table {
		for i, cell := range row {
			width := len(cell)
			if width > widths[i] {
				widths[i] = width
			}
		}
	}

	for _, row := range table {
		for i, cell := range row {
			// do not pad last line
			if len(row)-1 == i {
				fmt.Print(cell)
				continue
			}
			fmt.Printf("%-"+strconv.Itoa(widths[i])+"s  ", cell)
		}
		fmt.Print("\n")
	}
}
