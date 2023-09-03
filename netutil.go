package qianmo

import (
	"net"
	"os"
	"strings"

	"golang.org/x/net/nettest"
)

// FindInterfaceByName returns the interface with the given name.
func FindInterfaceByName(name string) (*net.Interface, error) {
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

// FindInterfaceByIP returns the interface with the given IP address.
func FindInterfaceByIP(ip string) (*net.Interface, error) {
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

// FindInterfaceByMAC returns the interface with the given MAC address.
func FindInterfaceByMAC(mac net.HardwareAddr) (*net.Interface, error) {
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

// FindLoopbackInterface returns the loopback interface.
func FindLoopbackInterface() (*net.Interface, error) {
	return nettest.LoopbackInterface()
}

// FindMacByIP returns the MAC address of the interface with the given IP address.
func FindMacByIP(ip string) (net.HardwareAddr, error) {
	iface, err := FindInterfaceByIP(ip)
	if err != nil {
		return nil, err
	}

	return iface.HardwareAddr, nil
}

// FindMacByName returns the MAC address of the interface with the given name.
func FindMacByName(name string) (net.HardwareAddr, error) {
	iface, err := FindInterfaceByName(name)
	if err != nil {
		return nil, err
	}

	return iface.HardwareAddr, nil
}

// FindAddrs returns the IP addresses of the interface with the given iface 	name.
func FindAddrs(name string) []string {
	iface, err := FindInterfaceByName(name)
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

// FindNonLoopbackAddrs returns the non-loopback IP addresses of interfaces.
func FindNonLoopbackAddrs() []string {
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

// FindLoopbackAddrs returns the loopback IP addresses of interfaces.
func FindLoopbackAddrs() []string {
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

// FindHostIP returns the IP addresses of the host.
func FindHostIP() ([]string, error) {
	name, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	return net.LookupHost(name)
}

// FindHostFirstIPv6 returns the first non-loopback IPv4 address of the host.
func FindHostFirstIPv4() (string, error) {
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

// FindHostFirstIPv6 returns the first non-loopback IPv6 address of the host.
func FindHostFirstIPv6() (string, error) {
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
