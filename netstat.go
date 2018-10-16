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

// Netstat should point to a file in the proc filesystem where information
// about open sockets can be gathered.
type Netstat string

// Entry contains the gathered information.
type Entry struct {
	// Exe contains the path to the process.
	// Exe is empty if there was an error reading /proc/pid/exe.
	Exe string
	// Cmdline contains the complete command line for the process.
	// Returns an empty array if /proc/pid/cmdline can't be read.
	Cmdline []string
	// Pid contains the pid of the process. Is zero if open connection can't be assigned to a pid.
	Pid int

	// Inode contains the inode for the open connection.
	Inode uint64

	// IP holds the local IP for the connection.
	IP net.IP
	// Port holds the local port for the connection.
	Port int64
	// RemoteIP holds the remote IP for the connection.
	RemoteIP net.IP
	// RemotePort holds the remote port for the connection.
	RemotePort int64
}

// ProcRoot should point to the root of the proc file system
var ProcRoot = "/proc"

var (
	// TCP contains the standard location to read open TCP IPv4 connections.
	TCP = Netstat("net/tcp")
	// TCP6 contains the standard location to read open TCP IPv6 connections.
	TCP6 = Netstat("net/tcp6")
	// UDP contains the standard location to read open UDP IPv4 connections.
	UDP = Netstat("net/udp")
	// UDP6 contains the standard location to read open UDP IPv6 connections.
	UDP6 = Netstat("net/udp6")
)

var (
	procFdLinkParseType1 = regexp.MustCompile(`^socket:\[(\d+)\]$`)
	procFdLinkParseType2 = regexp.MustCompile(`^\[0000\]:(\d+)$`)
)

// Entries queries the given /proc/net file and returns the found entries.
// Returns an error if the /proc/net file can't be read.
func (n Netstat) Entries() ([]Entry, error) {
	inodeToPid := make(chan map[uint64]int)

	go func() {
		inodeToPid <- procFdInodeToPid()
	}()

	lines, err := n.readProcNetFile()
	if err != nil {
		return nil, err
	}

	entries := procNetToEntries(lines, <-inodeToPid)

	return entries, nil
}

func procNetToEntries(lines [][]string, inodeToPid map[uint64]int) []Entry {
	entries := make([]Entry, 0, len(lines))
	for _, line := range lines {
		localIPPort := strings.Split(line[1], ":")
		remoteIPPort := strings.Split(line[2], ":")
		inode := parseInode(line[9])
		pid := inodeToPid[inode]

		entry := Entry{
			Exe:        procGetExe(pid),
			Cmdline:    procGetCmdline(pid),
			Pid:        pid,
			Inode:      inode,
			IP:         parseIP(localIPPort[0]),
			Port:       hexToDec(localIPPort[1]),
			RemoteIP:   parseIP(remoteIPPort[0]),
			RemotePort: hexToDec(remoteIPPort[1]),
		}

		entries = append(entries, entry)
	}
	return entries
}

func parseInode(num string) uint64 {
	inode, _ := strconv.ParseUint(num, 10, 64)
	return inode
}

func (n Netstat) readProcNetFile() ([][]string, error) {
	var lines [][]string

	f, err := os.Open(filepath.Join(ProcRoot, string(n)))
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
		return nil, fmt.Errorf("can't read proc file: %s has no content", n)
	}
	// Remove header line
	return lines[1:], nil
}

func lineParts(line string) []string {
	parts := strings.Split(line, " ")
	filtered := parts[:0]
	for _, part := range parts {
		if len(part) > 0 {
			filtered = append(filtered, part)
		}
	}
	return filtered
}

func hexToDec(hex string) int64 {
	dec, _ := strconv.ParseInt(hex, 16, 64)
	return dec
}

func parseIP(ip string) net.IP {
	return net.IP(parseIPSegments(ip))
}

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

	paths, err := filepath.Glob(filepath.Join(ProcRoot, "[0-9]*/fd/[0-9]*"))
	if err != nil {
		return inodeToPid
	}

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
