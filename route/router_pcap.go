package route

import (
	"net"

	"github.com/google/gopacket/routing"
	"github.com/smallnest/qianmo"
)

// RouteByPcap returns the interface, gateway and preferred source IP address for the given destination.
func RouteByPcap(dst string) (iface *net.Interface, gateway, preferredSrc net.IP, macAddr string, err error) {
	router, err := routing.New()
	if err != nil {
		return nil, nil, nil, "", err
	}

	iface, gateway, preferredSrc, err = router.Route(net.ParseIP(dst))
	if err != nil {
		return nil, nil, nil, "", err
	}

	macAddr, err = GetMacAddr(dst, gateway.String())

	return iface, gateway, preferredSrc, macAddr, err
}

// RouteByPcapWithSrc returns the interface, gateway and preferred source IP address for the given destination and source.
func RouteByPcapWithSrc(src, dst string) (iface *net.Interface, gateway, preferredSrc net.IP, macAddr string, err error) {
	router, err := routing.New()
	if err != nil {
		return nil, nil, nil, "", err
	}

	ifa, err := qianmo.GetInterfaceByIP(src)
	if err != nil {
		return nil, nil, nil, "", err
	}

	iface, gateway, preferredSrc, err = router.RouteWithSrc(ifa.HardwareAddr, net.ParseIP(src), net.ParseIP(dst))

	if err != nil {
		return nil, nil, nil, "", err
	}

	macAddr, err = GetMacAddr(dst, gateway.String())

	return iface, gateway, preferredSrc, macAddr, err
}
