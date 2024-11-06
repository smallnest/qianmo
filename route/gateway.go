package route

import (
	"fmt"
	"net"
	"os/exec"
	"strings"
	"time"
)

// GetGatewayIP returns the gateway IP address for a given IP and network interface name.
func GetGatewayIP(ip, interfaceName string) (string, error) {
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

	// check each line for the interface name and the target IP
	for _, line := range lines {
		// Skip lines that don't contain the interface name
		if !strings.Contains(line, interfaceName) {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}

		// check if the line contains the default gateway
		if fields[0] == "default" {
			// check if the line contains the "via" keyword
			for i, field := range fields {
				if field == "via" && i+1 < len(fields) {
					defaultGateway = fields[i+1]
					break
				}
			}
			continue
		}

		// Parse the network IP and mask
		_, ipNet, err := net.ParseCIDR(fields[0])
		if err != nil {
			continue
		}

		// Check if the target IP is in the network range
		if ipNet.Contains(targetIP) {
			// check if the line contains the "via" keyword, and return the gateway IP
			for i, field := range fields {
				if field == "via" && i+1 < len(fields) {
					return fields[i+1], nil
				}
			}
		}
	}

	// Return the default gateway if no gateway was found for the target IP
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

// GetMacAddr returns the MAC address for a given destination and gateway IP address.
func GetMacAddr(dst, gateway string) (string, error) {
	macAddr, err := GetMacAddrFromArpCache(gateway)
	if err != nil && strings.Contains(err.Error(), "not found") {
		conn, _ := net.DialTimeout("ip:icmp", dst, time.Second)
		if conn != nil {
			conn.Close()
		}
		macAddr, err = GetMacAddrFromArpCache(gateway)
	}
	return macAddr, err
}
