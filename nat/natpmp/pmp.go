package natpmp

import (
	"fmt"
	"log"
	"time"

	"github.com/jackpal/gateway"
	natpmp "github.com/jackpal/go-nat-pmp"
)

type NatPMP struct {
	c *natpmp.Client
}

func New() (*NatPMP, error) {
	np := new(NatPMP)

	gatewayIP, err := gateway.DiscoverGateway()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	np.c = natpmp.NewClient(gatewayIP)
	if _, err := np.c.GetExternalAddress(); err != nil {
		log.Println(err)
		return nil, err
	}

	return np, nil
}

func (np *NatPMP) IP() (string, error) {
	response, err := np.c.GetExternalAddress()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d.%d.%d.%d", response.ExternalIPAddress[0], response.ExternalIPAddress[1], response.ExternalIPAddress[2], response.ExternalIPAddress[3]), nil
}

func (np *NatPMP) AddMapping(protocol string, inner, outer uint16, lifetime time.Duration) error {
	_, err := np.c.AddPortMapping(protocol, int(inner), int(outer), int(lifetime))
	if err != nil {
		return err
	}
	return nil
}
