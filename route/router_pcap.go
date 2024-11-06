package route

import (
	"net"

	"github.com/google/gopacket/routing"
	"github.com/smallnest/qianmo"
)

// RouteByPcap returns the interface, gateway and preferred source IP address for the given destination.
func RouteByPcap(dst string) (iface *net.Interface, gateway, preferredSrc net.IP, err error) {
	router, err := routing.New()
	if err != nil {
		return nil, nil, nil, err
	}

	return router.Route(net.ParseIP(dst))
}

// RouteByPcapWithSrc returns the interface, gateway and preferred source IP address for the given destination and source.
func RouteByPcapWithSrc(src, dst string) (iface *net.Interface, gateway, preferredSrc net.IP, err error) {
	router, err := routing.New()
	if err != nil {
		return nil, nil, nil, err
	}

	ifa, err := qianmo.GetInterfaceByIP(src)
	if err != nil {
		return nil, nil, nil, err
	}
	return router.RouteWithSrc(ifa.HardwareAddr, net.ParseIP(src), net.ParseIP(dst))
}
