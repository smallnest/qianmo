//go:build linux
// +build linux

package route

import (
	"fmt"
	"net"
	"time"
)

// Route returns the network interface, gateway, and preferred source IP for a given destination address.
func Route(dst string) (iface *net.Interface, gateway, preferredSrc net.IP, macAddr string, err error) {
	// resolve the destination IP address
	dstIP := net.ParseIP(dst)
	if dstIP == nil {
		return nil, nil, nil, "", fmt.Errorf("invalid destination IP address: %s", dst)
	}

	var src string
	conn, _ := net.DialTimeout("ip:icmp", dst, time.Second)
	if conn != nil {
		src = conn.LocalAddr().String()
		conn.Close()
	}
	if src == "" {
		return nil, nil, nil, "", fmt.Errorf("failed to determine source IP address for destination: %s", dst)
	}
	srcIP := net.ParseIP(src)

	// get the list of network interfaces
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, nil, nil, "", err
	}

	// iterate through the interfaces to find a suitable route
	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			var ipNet *net.IPNet
			switch v := addr.(type) {
			case *net.IPNet:
				ipNet = v
			case *net.IPAddr:
				ipNet = &net.IPNet{IP: v.IP, Mask: net.CIDRMask(32, 32)} // Assuming /32 for single IP addresses
			}

			if ipNet != nil && ipNet.Contains(srcIP) {
				// if the destination IP is in the same subnet, we can use this interface
				preferredSrc = ipNet.IP
				gatewayIP, err := GetGatewayIP(preferredSrc.To4().String(), iface.Name)
				if err != nil {
					return nil, nil, nil, "", err
				}
				gateway = net.ParseIP(gatewayIP)

				macAddr, err = GetMacAddr(dst, gateway.String())

				return &iface, gateway, preferredSrc, macAddr, err
			}
		}
	}

	return nil, nil, nil, "", fmt.Errorf("no suitable route found for destination IP: %s", dst)
}

// RouteWithSrc returns the network interface, gateway, and preferred source IP for a given source and destination address.
func RouteWithSrc(src, dst string) (iface *net.Interface, gateway, preferredSrc net.IP, macAddr string, err error) {
	// Parse source and destination IP addresses
	srcIP := net.ParseIP(src)
	dstIP := net.ParseIP(dst)
	if srcIP == nil || dstIP == nil {
		return nil, nil, nil, "", fmt.Errorf("invalid source or destination IP address: src=%s, dst=%s", src, dst)
	}

	// Get the list of network interfaces
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, nil, nil, "", err
	}

	// Iterate through the interfaces to find a suitable route
	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		var ipNet *net.IPNet
		preferredSrc = nil

		// Check if the source IP is associated with the interface
		for _, addr := range addrs {
			switch v := addr.(type) {
			case *net.IPNet:
				ipNet = v
			case *net.IPAddr:
				ipNet = &net.IPNet{IP: v.IP, Mask: net.CIDRMask(32, 32)} // Assuming /32 for single IP addresses
			}

			// Check if the source IP is in the interface's subnet
			if ipNet != nil && ipNet.Contains(srcIP) {
				preferredSrc = ipNet.IP
				break
			}
		}

		// If we found a preferred source, check the destination IP
		if preferredSrc != nil && ipNet != nil && ipNet.Contains(dstIP) {
			// Get gateway IP
			gatewayIP, err := GetGatewayIP(src, iface.Name)
			if err != nil {
				return nil, nil, nil, "", err
			}

			gateway = net.ParseIP(gatewayIP)

			macAddr, err = GetMacAddr(dst, gateway.String())

			return &iface, gateway, preferredSrc, macAddr, err
		}
	}

	return nil, nil, nil, "", fmt.Errorf("no suitable route found for src: %s and dst: %s", src, dst)
}
