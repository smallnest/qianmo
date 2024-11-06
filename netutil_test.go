package qianmo

import (
	"net"
	"testing"

	"github.com/gookit/goutil/testutil/assert"
	"github.com/stretchr/testify/require"
)

func TestGetInterfaceByIP(t *testing.T) {
	iface, err := GetInterfaceByIP("127.0.0.1")
	require.NoError(t, err)
	assert.NotEmpty(t, iface)
}

func TestGetInterfaceByName(t *testing.T) {
	// 先获取一个有效的接口名称
	interfaces, err := net.Interfaces()
	require.NoError(t, err)
	require.NotEmpty(t, interfaces)

	iface, err := GetInterfaceByName(interfaces[0].Name)
	require.NoError(t, err)
	assert.NotEmpty(t, iface)
}

func TestGetNonLoopbackAddrs(t *testing.T) {
	addrs, err := GetNonLoopbackAddrs()
	require.NoError(t, err)
	assert.NotEmpty(t, addrs)
	t.Logf("addrs: %v", addrs)
}

func TestGetHostIP(t *testing.T) {
	ip, err := GetHostIP()
	require.NoError(t, err)
	assert.NotEmpty(t, ip)
	t.Logf("ip: %v", ip)
}

func TestGetFreeTCPPort(t *testing.T) {
	port, err := GetFreeTCPPort()
	require.NoError(t, err)

	assert.Gt(t, port, 0)
	assert.Lt(t, port, 65536)
}

func TestGetFreeUDPPort(t *testing.T) {
	port, err := GetFreeUDPPort()
	require.NoError(t, err)

	assert.Gt(t, port, 0)
	assert.Lt(t, port, 65536)
}
