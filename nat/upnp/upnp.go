package upnp

import (
	"context"
	"errors"
	"log"

	"github.com/snowmerak/twisted-lyfes/net/port"
	"gitlab.com/NebulousLabs/go-upnp"
)

type UPnP struct {
	d *upnp.IGD
	p uint16
}

func New(ctx context.Context) (*UPnP, error) {
	u := new(UPnP)

	p := port.MIN

	d, err := upnp.DiscoverCtx(ctx)
	if err != nil {
		return nil, err
	}

	for {
		if err := d.Forward(p, "peer2peer"); err != nil {
			p++
			if p > port.MAX {
				return nil, errors.New("no avaliable port")
			}
			continue
		}
		if port.IsAvaliable(p) {
			break
		}
	}

	u.d = d
	u.p = p

	return u, nil
}

func (u *UPnP) Clear() error {
	if err := u.d.Clear(u.p); err != nil {
		return err
	}
	log.Printf("Cleared port %d\n", u.p)
	return nil
}

func (u *UPnP) Port() uint16 {
	return u.p
}
