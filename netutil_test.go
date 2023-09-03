package qianmo

import (
	"errors"
	"testing"

	"github.com/gookit/goutil/testutil/assert"
	"github.com/stretchr/testify/require"
)

func TestInterfaceByMAC(t *testing.T) {
	mac, err := MacByIP("127.0.0.1")
	require.NoError(t, err)

	iface, err := InterfaceByMAC(mac)
	require.NoError(t, err)
	assert.NotEmpty(t, iface)
	assert.Equal(t, mac.String(), iface.HardwareAddr.String())

}

func TestMacByName(t *testing.T) {
	lookback, err := LoopbackInterface()
	require.NoError(t, err)

	mac, err := MacByName(lookback.Name)

	require.NoError(t, err)

	iface, err := InterfaceByMAC(mac)
	require.NoError(t, err)
	assert.NotEmpty(t, iface)
	assert.Equal(t, mac.String(), iface.HardwareAddr.String())
}

func TestInterfaceByName(t *testing.T) {
	lookback, err := LoopbackInterface()
	require.NoError(t, err)

	iface, err := InterfaceByName(lookback.Name)
	require.NoError(t, err)
	assert.NotEmpty(t, iface)
}

func TestAddrs(t *testing.T) {
	lookback, err := LoopbackInterface()
	require.NoError(t, err)

	addrs := Addrs(lookback.Name)
	assert.Gt(t, len(addrs), 0)
}

func TestNonLoopbackAddrs(t *testing.T) {
	addrs := NonLoopbackAddrs()
	assert.NotEmpty(t, addrs)
	t.Logf("addrs: %v", addrs)
}

func TestLoopbackAddrs(t *testing.T) {
	addrs := LoopbackAddrs()
	assert.NotEmpty(t, addrs)
	t.Logf("addrs: %v", addrs)
}

func TestHostIP(t *testing.T) {
	addrs, err := HostIP()
	require.NoError(t, err)
	assert.NotEmpty(t, addrs)
	t.Logf("addrs: %v", addrs)
}

func TestHostFirstIPv4(t *testing.T) {
	ip, err := HostFirstIPv4()
	require.NoError(t, err)
	assert.NotEmpty(t, ip)
	t.Logf("ip: %v", ip)
}

func TestHostFirstIPv6(t *testing.T) {
	ip, err := HostFirstIPv6()
	if errors.Is(err, ErrNotFound) {
		t.Skip(err)
	}
	require.NoError(t, err)
	assert.NotEmpty(t, ip)
	t.Logf("ip: %v", ip)
}
