package utils

import (
	"fmt"
	"net"
)

// GetInternalIPv4 返回本机的内网 IPv4 地址
func GetInternalIPv4() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, i := range interfaces {
		address, err := i.Addrs()
		if err != nil {
			return "", err
		}

		for _, addr := range address {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			// 排除回环地址
			if ip == nil || ip.IsLoopback() {
				continue
			}

			// 检查是否为 IPv4 地址
			if ip.To4() != nil {
				return ip.String(), nil
			}
		}
	}

	return "", fmt.Errorf("no internal IPv4 address found")
}
