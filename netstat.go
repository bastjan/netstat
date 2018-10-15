package netstat

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type Netstat string

type Entry struct {
	Name string

	Inode uint64

	IP         string
	Port       int64
	RemoteIP   string
	RemotePort int64
}

var (
	TCP  = Netstat("/proc/net/tcp")
	TCP6 = Netstat("/proc/net/tcp6")
	UDP  = Netstat("/proc/net/udp")
	UDP6 = Netstat("/proc/net/udp6")
)

func (n Netstat) Entries() ([]Entry, error) {
	lines, err := n.readProcFile()
	if err != nil {
		return nil, err
	}
	entries := make([]Entry, 0, len(lines))
	for _, line := range lines {
		localIPPort := strings.Split(line[1], ":")
		remoteIPPort := strings.Split(line[2], ":")

		entry := Entry{
			Inode:      parseInode(line[9]),
			IP:         parseIP(localIPPort[0]),
			Port:       hexToDec(localIPPort[1]),
			RemoteIP:   parseIP(localIPPort[0]),
			RemotePort: hexToDec(remoteIPPort[1]),
		}

		entries = append(entries, entry)
	}

	return entries, nil
}

func parseInode(num string) uint64 {
	inode, _ := strconv.ParseUint(num, 10, 64)
	return inode
}

func (n Netstat) readProcFile() ([][]string, error) {
	var lines [][]string

	f, err := os.Open(string(n))
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
		return nil, errors.New("can't read proc file: file has no contents")
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

func parseIP(ip string) string {
	switch len(ip) {
	case 8:
		return parseIP4(ip)
	case 32:
		return parseIP6(ip)
	default:
		return ""
	}
}

func parseIP4(ip string) string {
	seg := parseIPSegments(ip)
	return fmt.Sprintf("%d.%d.%d.%d", seg[0], seg[1], seg[2], seg[3])
}

func parseIP6(ip string) string {
	seg := parseIPSegments(ip)
	return fmt.Sprintf("%x%x:%x%x:%x%x:%x%x:%x%x:%x%x:%x%x:%x%x",
		seg[0], seg[1], seg[2], seg[3],
		seg[4], seg[5], seg[6], seg[7],
		seg[8], seg[9], seg[10], seg[11],
		seg[12], seg[13], seg[14], seg[15])
}

func parseIPSegments(ip string) []uint8 {
	segments := make([]uint8, 0, len(ip)/2)
	for i := len(ip); i > 0; i -= 2 {
		seg, _ := strconv.ParseUint(ip[i-2:i], 16, 8)
		segments = append(segments, uint8(seg))
	}
	return segments
}
