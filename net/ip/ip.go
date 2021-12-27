package ip

import (
	"fmt"
	"math"
	"net"
	"strconv"
	"strings"
)

func GetBitFrom(ip *net.IP) (uint32, error) {
	if ip == nil {
		return 0, fmt.Errorf("ip.GetBitFrom: ip is nil")
	}
	ipm := uint32(0)
	for i, s := range strings.Split(ip.String(), ".") {
		v, err := strconv.Atoi(s)
		if err != nil {
			return 0, fmt.Errorf("ip.GetBitFrom: %w", err)
		}
		ipm |= uint32(v) << uint32(24-i*8)
	}
	return ipm, nil
}

func GetMaskBitFrom(ip *net.IP) uint32 {
	maskBitSize, _ := ip.DefaultMask().Size()
	maskBit := uint32(0)
	for i := 0; i < maskBitSize; i++ {
		maskBit |= 1 << uint32(31-i)
	}
	return maskBit
}

func GetFirstIP(ip *net.IP) (*net.IP, error) {
	if ip == nil {
		return nil, fmt.Errorf("ip.GetFirstIp: ip is nil")
	}
	ipm, err := GetBitFrom(ip)
	if err != nil {
		return nil, fmt.Errorf("ip.GetFirstIp: %w", err)
	}
	maskBit := GetMaskBitFrom(ip)
	ipm &= maskBit
	p := net.IPv4(byte(ipm>>24), byte(ipm>>16), byte(ipm>>8), byte(ipm))
	n, err := GetNextIP(&p)
	return n, nil
}

func GetNextIP(ip *net.IP) (*net.IP, error) {
	ipm, err := GetBitFrom(ip)
	if err != nil {
		return nil, fmt.Errorf("ip.GetNextIP: %w", err)
	}
	maskBit := GetMaskBitFrom(ip)
	ipm += 1
	if (ipm & (maskBit ^ math.MaxUint32)) == 0 {
		return nil, fmt.Errorf("ip.GetNextIP: no next ip")
	}
	p := net.IPv4(byte(ipm>>24), byte(ipm>>16), byte(ipm>>8), byte(ipm))
	return &p, nil
}

func GetLocalIPs() ([]*net.IP, error) {
	ips := []*net.IP{}
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("ip.GetLocalIPs: %w", err)
	}
	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, fmt.Errorf("ip.GetLocalIPs: %w", err)
		}
		for _, addr := range addrs {
			ip, _, err := net.ParseCIDR(addr.String())
			if err != nil {
				return nil, fmt.Errorf("ip.GetLocalIPs: %w", err)
			}
			if ip.To4() != nil && !ip.IsLoopback() {
				ips = append(ips, &ip)
			}
		}
	}
	return ips, nil
}
