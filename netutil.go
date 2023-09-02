package qianmo

import (
	"net"
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
