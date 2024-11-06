package route

import (
	"fmt"
	"os/exec"
	"strings"
)

// GetGatewayIP returns the gateway IP address for a given network interface name.
func GetGatewayIP(interfaceName string) (string, error) {
	// Use `ip route` command to get the gateway IP
	cmd := exec.Command("ip", "route")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, interfaceName) {
			fields := strings.Fields(line)
			for i, field := range fields {
				if field == "default" && i+2 < len(fields) {
					return fields[i+2], nil // Gateway IP is the second field after "default"
				}
			}
		}
	}
	return "", fmt.Errorf("gateway not found for interface: %s", interfaceName)
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
