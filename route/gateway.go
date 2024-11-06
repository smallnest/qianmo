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
