package qianmo

import (
	"net"
	"os"
	"strings"

	"golang.org/x/net/nettest"
)

// InterfaceByName returns the interface with the given name.
func InterfaceByName(name string) (*net.Interface, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, iface := range interfaces {
		if strings.EqualFold(iface.Name, name) {
			return &iface, nil
		}
	}

	return nil, ErrNotFound
}

// InterfaceByIP returns the interface with the given IP address.
func InterfaceByIP(ip string) (*net.Interface, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, iface := range interfaces {
		if addrs, err := iface.Addrs(); err == nil {

			for _, addr := range addrs {

				if ipNet, ok := addr.(*net.IPNet); ok {
					if ipNet.IP.String() == ip {
						return &iface, nil
					}
				}
			}
		}
	}

	return nil, ErrNotFound
}

// InterfaceByMAC returns the interface with the given MAC address.
func InterfaceByMAC(mac net.HardwareAddr) (*net.Interface, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, iface := range interfaces {
		if strings.EqualFold(iface.HardwareAddr.String(), mac.String()) {
			return &iface, nil
		}
	}

	return nil, ErrNotFound
}

// LoopbackInterface returns the loopback interface.
func LoopbackInterface() (*net.Interface, error) {
	return nettest.LoopbackInterface()
}

// MacByIP returns the MAC address of the interface with the given IP address.
func MacByIP(ip string) (net.HardwareAddr, error) {
	iface, err := InterfaceByIP(ip)
	if err != nil {
		return nil, err
	}

	return iface.HardwareAddr, nil
}

// MacByName returns the MAC address of the interface with the given name.
func MacByName(name string) (net.HardwareAddr, error) {
	iface, err := InterfaceByName(name)
	if err != nil {
		return nil, err
	}

	return iface.HardwareAddr, nil
}

// Addrs returns the IP addresses of the interface with the given iface 	name.
func Addrs(name string) []string {
	iface, err := InterfaceByName(name)
	if err != nil {
		return nil
	}

	addrs, err := iface.Addrs()
	if err != nil {
		return nil
	}

	var ips []string
	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok {
			ips = append(ips, ipNet.IP.String())
		}
	}

	return ips
}

// NonLoopbackAddrs returns the non-loopback IP addresses of interfaces.
func NonLoopbackAddrs() []string {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil
	}

	var ips []string
	for _, iface := range interfaces {
		if iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			if ipNet, ok := addr.(*net.IPNet); ok {
				ips = append(ips, ipNet.IP.String())
			}
		}
	}

	return ips
}

// LoopbackAddrs returns the loopback IP addresses of interfaces.
func LoopbackAddrs() []string {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil
	}

	var ips []string
	for _, iface := range interfaces {
		if iface.Flags&net.FlagLoopback == 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			if ipNet, ok := addr.(*net.IPNet); ok {
				ips = append(ips, ipNet.IP.String())
			}
		}
	}

	return ips
}

// HostIP returns the IP addresses of the host.
func HostIP() ([]string, error) {
	name, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	return net.LookupHost(name)
}

// HostFirstIPv6 returns the first non-loopback IPv4 address of the host.
func HostFirstIPv4() (string, error) {
	name, err := os.Hostname()
	if err != nil {
		return "", err
	}

	addrs, err := net.LookupHost(name)
	if err != nil {
		return "", err
	}

	for _, addr := range addrs {
		ip := net.ParseIP(addr)
		if ip.To4() != nil && !ip.IsLoopback() {
			return addr, nil
		}
	}

	return "", ErrNotFound
}

// HostFirstIPv6 returns the first non-loopback IPv6 address of the host.
func HostFirstIPv6() (string, error) {
	name, err := os.Hostname()
	if err != nil {
		return "", err
	}

	addrs, err := net.LookupHost(name)
	if err != nil {
		return "", err
	}

	for _, addr := range addrs {
		ip := net.ParseIP(addr)

		if ip.To4() == nil && !ip.IsLoopback() && !strings.HasPrefix(addr, "fe80:") {
			return addr, nil
		}
	}

	return "", ErrNotFound
}
