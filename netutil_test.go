package qianmo

import (
	"net"
	"testing"

	"github.com/gookit/goutil/testutil/assert"
	"github.com/stretchr/testify/require"
)

func TestFindInterfaceByMAC(t *testing.T) {
	mac, err := FindMacByIP("127.0.0.1")
	require.NoError(t, err)

	iface, err := FindInterfaceByMAC(mac)
	require.NoError(t, err)
	assert.NotEmpty(t, iface)
	assert.Equal(t, mac.String(), iface.HardwareAddr.String())

}

func TestFindMacByName(t *testing.T) {
	lookback, err := FindLoopbackInterface()
	require.NoError(t, err)

	mac, err := FindMacByName(lookback.Name)

	require.NoError(t, err)

	iface, err := FindInterfaceByMAC(mac)
	require.NoError(t, err)
	assert.NotEmpty(t, iface)
	assert.Equal(t, mac.String(), iface.HardwareAddr.String())
}

func TestFindInterfaceByName(t *testing.T) {
	lookback, err := FindLoopbackInterface()
	require.NoError(t, err)

	iface, err := FindInterfaceByName(lookback.Name)
	require.NoError(t, err)
	assert.NotEmpty(t, iface)
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

func TestFindAddrs(t *testing.T) {
	lookback, err := FindLoopbackInterface()
	require.NoError(t, err)

	addrs := FindAddrs(lookback.Name)
	assert.Gt(t, len(addrs), 0)
}
