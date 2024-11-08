package main

import (
	"flag"
	"fmt"

	"github.com/smallnest/qianmo"
	"github.com/smallnest/qianmo/route"
)

func panicErr(err error) {
	if err != nil {
		panic(err)
	}
}

var (
	dstIP = flag.String("d", "8.8.8.8", "destination IP address")
)

func main() {
	flag.Parse()

	// interfaces
	{
		ifaces, err := qianmo.GetAllInterfaces()
		panicErr(err)
		fmt.Printf("interfaces: %v\n", ifaces)

		name := ifaces[0]
		if len(ifaces) > 1 {
			name = ifaces[1]
		}

		iface, err := qianmo.GetInterfaceByName(name)
		panicErr(err)
		fmt.Printf("interface by name: %v\n", iface)

		ip, err := qianmo.GetHostIP()
		panicErr(err)
		fmt.Printf("host ip: %v\n", ip)

		iface, err = qianmo.GetInterfaceByIP(ip)
		panicErr(err)
		fmt.Printf("interface by ip: %v\n", iface)
	}

	// addrs
	{
		addrs, err := qianmo.GetAddrs()
		panicErr(err)
		fmt.Printf("addrs: %v\n", addrs)

		addrs, err = qianmo.GetNonLoopbackAddrs()
		panicErr(err)
		fmt.Printf("non loopback addrs: %v\n", addrs)

		tcpPort, err := qianmo.GetFreeTCPPort()
		panicErr(err)
		fmt.Printf("free tcp port: %v\n", tcpPort)

		udpPort, err := qianmo.GetFreeUDPPort()
		panicErr(err)
		fmt.Printf("free udp port: %v\n", udpPort)
	}

	// route
	{
		iface, gateway, localIP, macAddr, err := route.Route(*dstIP)
		panicErr(err)

		fmt.Printf("route: iface=%v, gateway=%v, localIP=%v, macAddr=%v\n", iface, gateway, localIP, macAddr)

		src, err := qianmo.GetHostIP()
		panicErr(err)
		iface, gateway, localIP, macAddr, err = route.RouteWithSrc(src, *dstIP)
		panicErr(err)
		fmt.Printf("route with src: iface=%v, gateway=%v, localIP=%v, macAddr=%v\n", iface, gateway, localIP, macAddr)
	}
}
