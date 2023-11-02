// Package irtc provides a function to convert IP ranges to CIDR list.
package irtc

import (
	"net"
	"strconv"
	"strings"
)

// ConvertIPRangeToCIDRs converts IP range to CIDRs.
//
// It returns the CIDRs.
// For example, ConvertIPRangeToCIDRs("10.0.0.5", "10.0.0.7") returns
// the CIDRs 10.0.0.5/32 and 10.0.0.6/31.
func ConvertIPRangeToCIDRs(begin, end string) ([]string, error) {
	beginIP := net.ParseIP(begin)
	if beginIP == nil {
		return nil, &net.ParseError{Type: "IP address", Text: begin}
	}

	endIP := net.ParseIP(end)
	if endIP == nil {
		return nil, &net.ParseError{Type: "IP address", Text: end}
	}

	if !lowerEqual(beginIP, endIP) {
		return nil, nil
	}

	if strings.Contains(begin, ".") {
		// IPv4

		// No matter IPv4 or IPv6, ParseIP() returns a 16-byte slice, but
		// only the last 4 bytes are needed while processing IPv4 address.
		beginIP = beginIP[12:]
		endIP = endIP[12:]
	}

	return convertIPRangeToCIDRs(beginIP, endIP)
}

func convertIPRangeToCIDRs(begin, end net.IP) ([]string, error) {
	index := -1
	for i := 0; i < len(begin); i++ {
		if begin[i]^end[i] == 0 {
			continue
		}
		index = i
		break
	}

	if index == -1 {
		return []string{genCIDR(begin, 8*len(begin))}, nil
	}

	l := 0
	v := begin[index] ^ end[index]
	for v > 0 {
		l++
		v >>= 1
	}

	firstDiffPos := 8*index + 8 - l

	ip := make(net.IP, len(begin))
	copy(ip, begin)
	minAddr := genMinAddress(ip, firstDiffPos)
	maxAddr := genMaxAddress(ip, firstDiffPos)
	if lowerEqual(begin, minAddr) && lowerEqual(maxAddr, end) {
		return []string{genCIDR(ip, firstDiffPos)}, nil
	}

	var r []string
	dfs(begin, end, ip, firstDiffPos, &r)

	return r, nil
}

func dfs(begin, end, ip net.IP, pos int, result *[]string) {
	if pos >= 8*len(ip) {
		return
	}

	setZero(ip, pos)
	minAddr := genMinAddress(ip, pos+1)
	maxAddr := genMaxAddress(ip, pos+1)
	if lowerEqual(begin, minAddr) && lowerEqual(maxAddr, end) {
		*result = append(*result, genCIDR(ip, pos+1))
	} else {
		dfs(begin, end, ip, pos+1, result)
	}

	setOne(ip, pos)
	minAddr = genMinAddress(ip, pos+1)
	maxAddr = genMaxAddress(ip, pos+1)
	if lowerEqual(begin, minAddr) && lowerEqual(maxAddr, end) {
		*result = append(*result, genCIDR(ip, pos+1))
		return
	}
	dfs(begin, end, ip, pos+1, result)
}

func genCIDR(netIP net.IP, totalPrefix int) string {
	if totalPrefix == 8*len(netIP) {
		return netIP.String() + "/" + strconv.Itoa(totalPrefix)
	}

	index := totalPrefix >> 3
	prefix := totalPrefix % 8

	ip := make(net.IP, len(netIP))
	copy(ip, netIP)
	ip[index] = (ip[index] >> (8 - prefix)) << (8 - prefix)

	for i := index + 1; i < len(ip); i++ {
		ip[i] = 0x00
	}

	return ip.String() + "/" + strconv.Itoa(totalPrefix)
}

func genMinAddress(ip net.IP, pos int) net.IP {
	minAddr := make(net.IP, len(ip))
	copy(minAddr, ip)

	if pos >= 8*len(minAddr) {
		return minAddr
	}

	index := pos >> 3
	miniPos := pos % 8

	mask := byte(0xff) << (8 - miniPos)
	minAddr[index] &= mask

	for i := index + 1; i < len(minAddr); i++ {
		minAddr[i] = 0x00
	}

	return minAddr
}

func genMaxAddress(ip net.IP, pos int) net.IP {
	maxAddr := make(net.IP, len(ip))
	copy(maxAddr, ip)

	if pos >= 8*len(maxAddr) {
		return maxAddr
	}

	index := pos >> 3
	miniPos := pos % 8

	mask := byte(0xff) >> miniPos
	maxAddr[index] |= mask

	for i := index + 1; i < len(maxAddr); i++ {
		maxAddr[i] = 0xff
	}

	return maxAddr
}

func lowerEqual(a net.IP, b net.IP) bool {
	for i := 0; i < len(a); i++ {
		if a[i] == b[i] {
			continue
		}

		if a[i] > b[i] {
			return false
		}

		if a[i] < b[i] {
			return true
		}
	}

	return true
}

func setOne(ip net.IP, pos int) {
	index := pos >> 3
	miniPos := pos % 8
	ip[index] |= byte(0xff) >> 7 << (8 - miniPos - 1)
}

func setZero(ip net.IP, pos int) {
	index := pos >> 3
	miniPos := pos % 8
	ip[index] &= byte(0xff) << (8 - miniPos)
}
