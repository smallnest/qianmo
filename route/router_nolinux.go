//go:build !linux
// +build !linux

package route

import (
	"net"
)

// Route returns the interface, gateway and preferred source IP address for the given destination.
func Route(dst string) (iface *net.Interface, gateway, preferredSrc net.IP, macAddr string, err error) {
	panic("Route only implemented in linux")
}

// RouteWithSrc returns the interface, gateway and preferred source IP address for the given destination and source.
func RouteWithSrc(src, dst string) (iface *net.Interface, gateway, preferredSrc net.IP, macAddr string, err error) {
	panic("RouteWithSrc only implemented in linux")
}
