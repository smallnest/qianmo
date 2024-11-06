package route

import (
	"fmt"
	"net"
	"os/exec"
	"strings"
)

// GetGatewayIP returns the gateway IP address for a given IP and network interface name.
func GetGatewayIP(ip, interfaceName string) (string, error) {
	// 解析目标 IP
	targetIP := net.ParseIP(ip)
	if targetIP == nil {
		return "", fmt.Errorf("invalid IP address: %s", ip)
	}

	// Use `ip route` command to get the gateway IP
	cmd := exec.Command("ip", "route")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(output), "\n")
	var defaultGateway string

	// 首先查找特定网段的路由
	for _, line := range lines {
		// 检查行是否包含指定的接口名
		if !strings.Contains(line, interfaceName) {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}

		// 检查是否是默认路由
		if fields[0] == "default" {
			// 保存默认网关
			for i, field := range fields {
				if field == "via" && i+1 < len(fields) {
					defaultGateway = fields[i+1]
					break
				}
			}
			continue
		}

		// 解析网段
		_, ipNet, err := net.ParseCIDR(fields[0])
		if err != nil {
			continue
		}

		// 检查目标 IP 是否在这个网段内
		if ipNet.Contains(targetIP) {
			// 检查是否包含 "via" 关键字，它后面的字段就是网关IP
			for i, field := range fields {
				if field == "via" && i+1 < len(fields) {
					return fields[i+1], nil
				}
			}
		}
	}

	// 如果找到了默认网关，则返回默认网关
	if defaultGateway != "" {
		return defaultGateway, nil
	}

	return "", fmt.Errorf("gateway not found for IP %s on interface: %s", ip, interfaceName)
}

// GetMacAddrFromArpCache returns the MAC address for a given IP address.
func GetMacAddrFromArpCache(ip string) (string, error) {
	cmd := exec.Command("arp", "-n", ip)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, ip) {
			fields := strings.Fields(line)
			if len(fields) > 2 {
				return fields[2], nil // MAC address is usually the third field
			}
		}
	}
	return "", fmt.Errorf("MAC address not found for IP: %s", ip)
}

// Use Route
//
// // getGatewayIPAndMAC returns the gateway IP and MAC address for a given remote IP address.
// func getGatewayIPAndMAC(remoteIP string) (string, string, error) {
// 	// Get the list of network interfaces
// 	interfaces, err := net.Interfaces()
// 	if err != nil {
// 		return "", "", err
// 	}

// 	// Resolve the remote IP address
// 	ip := net.ParseIP(remoteIP)
// 	if ip == nil {
// 		return "", "", fmt.Errorf("invalid remote IP address: %s", remoteIP)
// 	}

// 	// Iterate through the interfaces to find the correct gateway
// 	for _, iface := range interfaces {
// 		// Filter out non-up interfaces and loopback interfaces
// 		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
// 			continue
// 		}

// 		// Get the interface's addresses
// 		addrs, err := iface.Addrs()
// 		if err != nil {
// 			continue
// 		}

// 		for _, addr := range addrs {
// 			// Only handle IPv4 addresses
// 			if ipNet, ok := addr.(*net.IPNet); ok && ipNet.IP.To4() != nil {
// 				// Check if the remote IP is in the same subnet as the interface
// 				if ipNet.Contains(ip) || ipNet.IP.Equal(ip) {
// 					// If the remote IP is in the same subnet, get its MAC address
// 					macAddr, err := GetMacAddress(remoteIP)
// 					if err == nil {
// 						return "", macAddr, nil // Return MAC address if found
// 					}

// 					// If not, get the gateway IP for this interface
// 					gatewayIP, err := GetGatewayIP(iface.Name)
// 					if err != nil {
// 						return "", "", err
// 					}

// 					macAddr = iface.HardwareAddr.String() // MAC address of the interface
// 					return gatewayIP, macAddr, nil
// 				}
// 			}
// 		}
// 	}
// 	return "", "", fmt.Errorf("no suitable gateway found for remote IP: %s", remoteIP)
// }
