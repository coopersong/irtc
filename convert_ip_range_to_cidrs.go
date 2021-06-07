package irtc

import (
    "fmt"
    "net"
    "strings"
)

// divide [ipBegin, ipEnd] to cidr list
func divideToCidrList(ipBegin, ipEnd []byte) ([]string, error) {
    ipBinaryBegin := convertIpToBoolSlice(ipBegin)
    ipBinaryEnd := convertIpToBoolSlice(ipEnd)
    startPos, err := firstDiffPos(ipBinaryBegin, ipBinaryEnd)
    if err != nil {
        // begin address and end address are the same
        prefixLength := len(ipBinaryBegin)
        cidr, err := genCidr(ipBinaryBegin, prefixLength)
        if err != nil {
            return nil, err
        }
        return []string{cidr}, nil
    }
    ip := copyBoolSlice(ipBinaryBegin)
    err = dfs(ip, startPos, true, ipBinaryBegin, ipBinaryEnd)
    if err != nil {
        return nil, err
    }
    return cidrList, nil
}

// convert ip range to cidrs
// for example, we can divide [10.0.0.5, 10.0.0.7] to 10.0.0.5/32 and 10.0.0.6/31
func ConvertIpRangeToCidrs(begin, end string) ([]string, error) {
    // reset cidrList
    cidrList = nil
    // The type of ipBegin is []byte
    ipBegin := net.ParseIP(begin)
    if ipBegin == nil {
        return nil, fmt.Errorf("invalid IP Begin format")
    }
    // The type of ipEnd is []byte
    ipEnd := net.ParseIP(end)
    if ipEnd == nil {
        return nil, fmt.Errorf("invalid IP End format")
    }
    if strings.Contains(begin, ".") {
        // IPv4

        // no matter IPv4 or IPv6, what we get from function net.ParseIP() is a 16-byte
        // byte slice, but we only need the last 4 bytes when we process IPv4 address.
        ipBegin = ipBegin[12:]
        ipEnd = ipEnd[12:]
    }
    cidrs, err := divideToCidrList(ipBegin, ipEnd)
    if err != nil {
        return nil, err
    }
    return cidrs, nil
}