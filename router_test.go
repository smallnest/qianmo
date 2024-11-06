package qianmo

import (
	"testing"

	"github.com/gookit/goutil/testutil/assert"
	"github.com/stretchr/testify/require"
)

func TestRoute(t *testing.T) {
	iface, gateway, ip, err := Route("114.114.114.114")
	require.NoError(t, err)
	assert.NotEmpty(t, iface, "iface: %s", iface)
	assert.NotEmpty(t, gateway, "gateway: %s", gateway)
	assert.NotEmpty(t, ip, "ip: %s", ip)

	localIP, err := GetHostIP()
	require.NoError(t, err)

	iface, gateway, ip, err = RouteWithSrc(localIP, "114.114.114.114")
	require.NoError(t, err)
	assert.NotEmpty(t, iface, "iface: %s", iface)
	assert.NotEmpty(t, gateway, "gateway: %s", gateway)
	assert.NotEmpty(t, ip, "ip: %s", ip)
}
