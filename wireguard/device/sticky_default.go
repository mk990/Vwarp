//go:build !linux

package device

import (
	"github.com/mk990/Vwarp/wireguard/conn"
	"github.com/mk990/Vwarp/wireguard/rwcancel"
)

func (device *Device) startRouteListener(bind conn.Bind) (*rwcancel.RWCancel, error) {
	return nil, nil
}
