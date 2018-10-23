/*
Package netstat helps you query open network connections.

Netstat searches the proc filesystem to gather information about open network connections and the
associated processes.

There is currently no support planned for Mac OS or *BSD without procfs support.
*/
package netstat

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// Protocol points to a file in the proc filesystem where information
// about open sockets can be gathered.
type Protocol struct {
	// Name contains the protocol name. Not used internally.
	Name string
	// RelPath is the proc file path relative to ProcRoot
	RelPath string
}

// ProcRoot should point to the root of the proc file system
var ProcRoot = "/proc"

var (
	// TCP contains the standard location to read open TCP IPv4 connections.
	TCP = &Protocol{"tcp", "net/tcp"}
	// TCP6 contains the standard location to read open TCP IPv6 connections.
	TCP6 = &Protocol{"tcp6", "net/tcp6"}
	// UDP contains the standard location to read open UDP IPv4 connections.
	UDP = &Protocol{"udp", "net/udp"}
	// UDP6 contains the standard location to read open UDP IPv6 connections.
	UDP6 = &Protocol{"udp6", "net/udp6"}
)

var (
	procFdLinkParseType1 = regexp.MustCompile(`^socket:\[(\d+)\]$`)
	procFdLinkParseType2 = regexp.MustCompile(`^\[0000\]:(\d+)$`)
)

// Connections queries the given /proc/net file and returns the found connections.
// Returns an error if the /proc/net file can't be read.
func (p *Protocol) Connections() ([]*Connection, error) {
	inodeToPid := make(chan map[uint64]int)

	go func() {
		inodeToPid <- procFdInodeToPid()
	}()

	lines, err := p.readProcNetFile()
	if err != nil {
		return nil, err
	}

	connections := p.procNetToConnections(lines, <-inodeToPid)

	return connections, nil
}

func (p *Protocol) procNetToConnections(lines [][]string, inodeToPid map[uint64]int) []*Connection {
	connections := make([]*Connection, 0, len(lines))
	for _, line := range lines {
		localIPPort := strings.Split(line[1], ":")
		remoteIPPort := strings.Split(line[2], ":")
		inode := parseUint64(line[9])
		pid := inodeToPid[inode]
		queues := strings.Split(line[4], ":")

		connection := &Connection{
			Exe:           procGetExe(pid),
			Cmdline:       procGetCmdline(pid),
			Pid:           pid,
			Inode:         inode,
			UserID:        line[7],
			IP:            parseIP(localIPPort[0]),
			Port:          parsePort(localIPPort[1]),
			RemoteIP:      parseIP(remoteIPPort[0]),
			RemotePort:    parsePort(remoteIPPort[1]),
			State:         tcpStatefromHex(line[3]),
			TransmitQueue: parseUint64(queues[0]),
			ReceiveQueue:  parseUint64(queues[1]),
			Protocol:      p,
		}

		connections = append(connections, connection)
	}
	return connections
}

func parseUint64(num string) uint64 {
	inode, _ := strconv.ParseUint(num, 10, 64)
	return inode
}

func (p *Protocol) readProcNetFile() ([][]string, error) {
	var lines [][]string

	path := filepath.Join(ProcRoot, p.RelPath)
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("can't open proc file: %s", err)
	}
	defer f.Close()

	r := bufio.NewReader(f)
	for {
		line, err := r.ReadBytes('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		lines = append(lines, lineParts(string(bytes.Trim(line, "\t\n "))))
	}
	if len(lines) == 0 {
		return nil, fmt.Errorf("can't read proc file: %s has no content", path)
	}
	// Remove header line
	return lines[1:], nil
}

// The values in a line are separated by one or more space.
// Split by space and remove all resulting empty strings.
// strings.Split("01   AB", " ") results in ["01", "", "", "AB"]
func lineParts(line string) []string {
	parts := strings.Split(line, " ")
	filtered := parts[:0]
	for _, part := range parts {
		if part != "" {
			filtered = append(filtered, part)
		}
	}
	return filtered
}

func parsePort(port string) int {
	// port number is an unsigned 16-bit integer
	dec, _ := strconv.ParseUint(port, 16, 16)
	return int(dec)
}

func parseIP(ip string) net.IP {
	return net.IP(parseIPSegments(ip))
}

// The IP address is encoded hexadecimal and in reverse order.
// Take two characters and parse then from back to front.
// 01 00 00 7F -> 127 0 0 1
func parseIPSegments(ip string) []uint8 {
	segments := make([]uint8, 0, len(ip)/2)
	for i := len(ip); i > 0; i -= 2 {
		seg, _ := strconv.ParseUint(ip[i-2:i], 16, 8)
		segments = append(segments, uint8(seg))
	}
	return segments
}

func procGetCmdline(pid int) []string {
	path := filepath.Join(ProcRoot, strconv.Itoa(pid), "cmdline")
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return []string{}
	}
	content = bytes.TrimRight(content, "\x00")
	return strings.Split(string(content), "\x00")
}

func procGetExe(pid int) string {
	path := filepath.Join(ProcRoot, strconv.Itoa(pid), "exe")
	target, err := os.Readlink(path)
	if err != nil {
		return ""
	}
	return target
}

func procFdInodeToPid() map[uint64]int {
	inodeToPid := make(map[uint64]int)

	// Ignoring error: The only possible error is bad pattern.
	paths, _ := filepath.Glob(filepath.Join(ProcRoot, "[0-9]*/fd/[0-9]*"))
	for _, link := range paths {
		target, err := os.Readlink(link)
		if err != nil {
			continue
		}

		pid := procFdExtractPid(link)
		inode, found := procFdExtractInode(target)
		if !found {
			continue
		}

		inodeToPid[inode] = pid
	}

	return inodeToPid
}

func procFdExtractPid(fdPath string) int {
	parts := strings.SplitN(fdPath, string(filepath.Separator), 4)
	pid, _ := strconv.ParseInt(parts[2], 10, 64)
	return int(pid)
}

func procFdExtractInode(fdLinkTarget string) (inode uint64, found bool) {
	match := procFdLinkParseType1.FindStringSubmatch(fdLinkTarget)
	if match == nil {
		match = procFdLinkParseType2.FindStringSubmatch(fdLinkTarget)
		if match == nil {
			return 0, false
		}
	}

	inode, _ = strconv.ParseUint(match[1], 10, 64)
	return inode, true
}
