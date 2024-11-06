//go:build !linux
// +build !linux

package qianmo

import (
	"net"
	"reflect"
	"testing"
)

func TestRouteWithSrc(t *testing.T) {
	type args struct {
		src string
		dst string
	}
	tests := []struct {
		name             string
		args             args
		wantIface        *net.Interface
		wantGateway      net.IP
		wantPreferredSrc net.IP
		wantErr          bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotIface, gotGateway, gotPreferredSrc, err := RouteWithSrc(tt.args.src, tt.args.dst)
			if (err != nil) != tt.wantErr {
				t.Errorf("RouteWithSrc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotIface, tt.wantIface) {
				t.Errorf("RouteWithSrc() gotIface = %v, want %v", gotIface, tt.wantIface)
			}
			if !reflect.DeepEqual(gotGateway, tt.wantGateway) {
				t.Errorf("RouteWithSrc() gotGateway = %v, want %v", gotGateway, tt.wantGateway)
			}
			if !reflect.DeepEqual(gotPreferredSrc, tt.wantPreferredSrc) {
				t.Errorf("RouteWithSrc() gotPreferredSrc = %v, want %v", gotPreferredSrc, tt.wantPreferredSrc)
			}
		})
	}
}
