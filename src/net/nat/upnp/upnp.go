package upnp

import (
	"context"
	"fmt"

	"gitlab.com/NebulousLabs/go-upnp"
)

type UPnP struct {
	d     *upnp.IGD
	Ports map[uint16]string
}

func New(ctx context.Context) (*UPnP, error) {
	u := new(UPnP)
	u.Ports = make(map[uint16]string)

	d, err := upnp.DiscoverCtx(ctx)
	if err != nil {
		return nil, fmt.Errorf("upnp.New: %w", err)
	}

	u.d = d

	return u, nil
}

func (u *UPnP) Forward(port uint16, desc string) error {
	if _, ok := u.Ports[port]; ok {
		return fmt.Errorf("upnp.Forward: port %d is already forwarded", port)
	}
	if err := u.d.Forward(port, desc); err != nil {
		return fmt.Errorf("upnp.Forward: %w", err)
	}
	u.Ports[port] = desc
	return nil
}

func (u *UPnP) Clear() error {
	for k := range u.Ports {
		if err := u.d.Clear(k); err != nil {
			return fmt.Errorf("upnp.Clear: %w", err)
		}
	}
	return nil
}

func (u *UPnP) ExternalIP() (string, error) {
	return u.d.ExternalIP()
}
