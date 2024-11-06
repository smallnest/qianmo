package qianmo

import (
	"net"
)

// Route returns the interface, gateway and source IP for reaching the target IP.
func Route(targetIP string) (*net.Interface, string, string, error) {
	if targetIP == "" {
		return nil, "", "", ErrInvalidParam
	}

	// Get all interfaces
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, "", "", err
	}

	// Try each interface
	for _, iface := range interfaces {
		// Skip loopback
		if iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			if ipNet, ok := addr.(*net.IPNet); ok {
				// Skip loopback and link-local addresses
				if ipNet.IP.IsLoopback() || ipNet.IP.IsLinkLocalUnicast() {
					continue
				}

				// Get gateway for this interface
				gateway, err := getDefaultGateway(iface.Name)
				if err != nil {
					continue
				}

				return &iface, gateway, ipNet.IP.String(), nil
			}
		}
	}

	return nil, "", "", ErrNotFound
}

// RouteWithSrc returns the interface, gateway and source IP for reaching the target IP
// using the specified source IP.
func RouteWithSrc(sourceIP, targetIP string) (*net.Interface, string, string, error) {
	if sourceIP == "" || targetIP == "" {
		return nil, "", "", ErrInvalidParam
	}

	// Get interface for source IP
	iface, err := GetInterfaceByIP(sourceIP)
	if err != nil {
		return nil, "", "", err
	}

	// Get gateway for this interface
	gateway, err := getDefaultGateway(iface.Name)
	if err != nil {
		return nil, "", "", err
	}

	return iface, gateway, sourceIP, nil
}

// getDefaultGateway returns the default gateway IP for the given interface.
// This is a placeholder implementation - you'll need to implement the actual
// gateway detection logic for your specific platform.
func getDefaultGateway(interfaceName string) (string, error) {
	// This is just a placeholder implementation
	// You should implement the actual gateway detection logic
	// based on your operating system:
	// - For Linux: parse /proc/net/route or use 'ip route' command
	// - For Windows: use GetIpForwardTable or similar Win32 API
	// - For macOS: use route -n get default or similar

	return "0.0.0.0", nil // Placeholder return
}
