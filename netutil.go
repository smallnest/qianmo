package qianmo

import (
	"net"
	"os"
	"strings"
)

// GetInterface returns the interface with the given name or IP address.
func GetAllInterfaces() ([]string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	var ifaces = make([]string, 0, len(interfaces))
	for _, iface := range interfaces {
		ifaces = append(ifaces, iface.Name)
	}

	return ifaces, nil
}

// GetInterfaceByName returns the interface with the given name.
func GetInterfaceByName(name string) (*net.Interface, error) {
	if name == "" {
		return nil, ErrInvalidParam
	}

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

// GetInterfaceByIP returns the interface with the given IP address.
func GetInterfaceByIP(ip string) (*net.Interface, error) {
	if ip == "" {
		return nil, ErrInvalidParam
	}

	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, iface := range interfaces {
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			if ipNet, ok := addr.(*net.IPNet); ok && ipNet.IP.String() == ip {
				return &iface, nil
			}
		}
	}

	return nil, ErrNotFound
}

// GetAddrs returns all IP addresses.
func GetAddrs() ([]string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	var ips []string
	for _, iface := range interfaces {
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}

		for _, addr := range addrs {
			if ipNet, ok := addr.(*net.IPNet); ok {
				ips = append(ips, ipNet.IP.String())
			}
		}
	}
	return ips, nil
}

// GetNonLoopbackAddrs returns all non-loopback IP addresses.
func GetNonLoopbackAddrs() ([]string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
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

	if len(ips) == 0 {
		return nil, ErrNotFound
	}
	return ips, nil
}

// GetHostIP returns the primary IP address of the host.
// It returns the first non-loopback IPv4 address if available, otherwise the first non-loopback IPv6 address.
func GetHostIP() (string, error) {
	name, err := os.Hostname()
	if err != nil {
		return "", err
	}

	addrs, err := net.LookupHost(name)
	if err != nil {
		return "", err
	}

	// 优先返回IPv4地址
	for _, addr := range addrs {
		ip := net.ParseIP(addr)
		if ip.To4() != nil && !ip.IsLoopback() {
			return addr, nil
		}
	}

	// 如果没有IPv4，返回第一个非回环IPv6地址
	for _, addr := range addrs {
		ip := net.ParseIP(addr)
		if ip.To4() == nil && !ip.IsLoopback() && !strings.HasPrefix(addr, "fe80:") {
			return addr, nil
		}
	}

	return "", ErrNotFound
}

// GetFreePort returns an available TCP port number.
func GetFreePort(proto string) (int, error) {
	addr, err := net.ResolveTCPAddr(proto, "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP(proto, addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()

	return l.Addr().(*net.TCPAddr).Port, nil
}
