//go:build linux
// +build linux

package route

import (
	"net"

	"github.com/google/gopacket/routing"
	"github.com/smallnest/qianmo"
)

// Route returns the interface, gateway and preferred source IP address for the given destination.
func Route(dst string) (iface *net.Interface, gateway, preferredSrc net.IP, err error) {
	router, err := routing.New()
	if err != nil {
		return nil, nil, nil, err
	}

	return router.Route(net.ParseIP(dst))
}

// RouteWithSrc returns the interface, gateway and preferred source IP address for the given destination and source.
func RouteWithSrc(src, dst string) (iface *net.Interface, gateway, preferredSrc net.IP, err error) {
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
