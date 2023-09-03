//go:build linux
// +build linux

package qianmo

import (
	"net"

	"github.com/google/gopacket/routing"
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
func RouteWithSrc(input net.HardwareAddr, src, dst string) (iface *net.Interface, gateway, preferredSrc net.IP, err error) {
	router, err := routing.New()
	if err != nil {
		return nil, nil, nil, err
	}

	return router.RouteWithSrc(input, net.ParseIP(src), net.ParseIP(dst))
}
