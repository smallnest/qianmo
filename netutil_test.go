package qianmo

import (
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

func TestFindAddrs(t *testing.T) {
	lookback, err := FindLoopbackInterface()
	require.NoError(t, err)

	addrs := FindAddrs(lookback.Name)
	assert.Gt(t, len(addrs), 0)
}

func TestFindNonLoopbackAddrs(t *testing.T) {
	addrs := FindNonLoopbackAddrs()
	assert.NotEmpty(t, addrs)
	t.Logf("addrs: %v", addrs)
}

func TestFindLoopbackAddrs(t *testing.T) {
	addrs := FindLoopbackAddrs()
	assert.NotEmpty(t, addrs)
	t.Logf("addrs: %v", addrs)
}

func TestFindHostIP(t *testing.T) {
	addrs, err := FindHostIP()
	require.NoError(t, err)
	assert.NotEmpty(t, addrs)
	t.Logf("addrs: %v", addrs)
}

func TestFindHostFirstIPv4(t *testing.T) {
	ip, err := FindHostFirstIPv4()
	require.NoError(t, err)
	assert.NotEmpty(t, ip)
	t.Logf("ip: %v", ip)
}

func TestFindHostFirstIPv6(t *testing.T) {
	ip, err := FindHostFirstIPv6()
	require.NoError(t, err)
	assert.NotEmpty(t, ip)
	t.Logf("ip: %v", ip)
}
