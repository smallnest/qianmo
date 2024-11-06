//go:build linux
// +build linux

package route

import (
	"testing"

	"github.com/gookit/goutil/testutil/assert"
	"github.com/smallnest/qianmo"
	"github.com/stretchr/testify/require"
)

func TestRoute(t *testing.T) {
	iface, gateway, ip, macAddr, err := Route("114.114.114.114")
	t.Logf("iface: %s, gateway: %s, ip: %s, macAddr: %s, err: %v", iface.Name, gateway, ip, macAddr, err)
	require.NoError(t, err)
	assert.NotEmpty(t, iface, "iface: %s", iface)
	assert.NotEmpty(t, gateway, "gateway: %s", gateway)
	assert.NotEmpty(t, ip, "ip: %s", ip)
	assert.NotEmpty(t, macAddr, "macAddr: %s", macAddr)

	localIP, err := qianmo.GetHostIP()
	require.NoError(t, err)

	iface, gateway, ip, macAddr, err = RouteByPcapWithSrc(localIP, "114.114.114.114")
	require.NoError(t, err)
	assert.NotEmpty(t, iface, "iface: %s", iface)
	assert.NotEmpty(t, gateway, "gateway: %s", gateway)
	assert.NotEmpty(t, ip, "ip: %s", ip)
	assert.NotEmpty(t, macAddr, "macAddr: %s", macAddr)
}

func TestRouteByPcap(t *testing.T) {
	iface, gateway, ip, macAddr, err := RouteByPcap("114.114.114.114")
	// t.Logf("iface: %s, gateway: %s, ip: %s, macAddr: %s, err: %v", iface.Name, gateway, ip, macAddr, err)
	require.NoError(t, err)
	assert.NotEmpty(t, iface, "iface: %s", iface)
	assert.NotEmpty(t, gateway, "gateway: %s", gateway)
	assert.NotEmpty(t, ip, "ip: %s", ip)
	assert.NotEmpty(t, macAddr, "macAddr: %s", macAddr)

	localIP, err := qianmo.GetHostIP()
	require.NoError(t, err)

	iface, gateway, ip, macAddr, err = RouteByPcapWithSrc(localIP, "114.114.114.114")
	require.NoError(t, err)
	assert.NotEmpty(t, iface, "iface: %s", iface)
	assert.NotEmpty(t, gateway, "gateway: %s", gateway)
	assert.NotEmpty(t, ip, "ip: %s", ip)
	assert.NotEmpty(t, macAddr, "macAddr: %s", macAddr)
}
