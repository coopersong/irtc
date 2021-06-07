package irtc

var cidrList []string

// depth-first search complete cidr list to represent ip range [ipBinaryBegin, ipBinaryEnd]
func dfs(ip []bool, pos int, isFirst bool, ipBinaryBegin, ipBinaryEnd []bool) error {
    if pos >= len(ipBinaryBegin) {
        return nil
    }

    var minAddress []bool
    var maxAddress []bool
    prefixLength := pos + 1
    if isFirst {
        minAddress = genMinAddress(ip, pos)
        maxAddress = genMaxAddress(ip, pos)
        if lowerEqual(ipBinaryBegin, minAddress) && lowerEqual(maxAddress, ipBinaryEnd) {
            cidr, err := genCidr(ip, prefixLength)
            if err != nil {
                return err
            }
            cidrList = append(cidrList, cidr)
            return nil
        }
        return dfs(ip, pos + 1, false, ipBinaryBegin, ipBinaryEnd)
    }

    ip[pos] = false
    minAddress = genMinAddress(ip, pos)
    maxAddress = genMaxAddress(ip, pos)
    if lowerEqual(ipBinaryBegin, minAddress) && lowerEqual(maxAddress, ipBinaryEnd) {
        cidr, err := genCidr(ip, prefixLength)
        if err != nil {
            return err
        }
        cidrList = append(cidrList, cidr)
    } else {
        err := dfs(ip, pos + 1, false, ipBinaryBegin, ipBinaryEnd)
        if err != nil {
            return nil
        }
    }

    ip[pos] = true
    minAddress = genMinAddress(ip, pos)
    maxAddress = genMaxAddress(ip, pos)
    if lowerEqual(ipBinaryBegin, minAddress) && lowerEqual(maxAddress, ipBinaryEnd) {
        cidr, err := genCidr(ip, prefixLength)
        if err != nil {
            return err
        }
        cidrList = append(cidrList, cidr)
    } else {
        err := dfs(ip, pos + 1, false, ipBinaryBegin, ipBinaryEnd)
        if err != nil {
            return err
        }
    }

    return nil
}