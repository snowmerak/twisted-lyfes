package discovery

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/snowmerak/generics-for-go/option"
	"github.com/snowmerak/twisted-lyfes/net/ip"
	"github.com/snowmerak/twisted-lyfes/net/port"
)

type Setting struct {
	Port      *option.Option[int]
	Limit     *option.Option[int]
	LimitTime *option.Option[time.Duration]
}

func ScanPort(ip net.IP, port int, limit time.Duration) (bool, error) {
	if ip == nil {
		return false, fmt.Errorf("discovery.ScanPort: ip is nil")
	}
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", ip.String(), port), limit)
	if err != nil {
		return false, fmt.Errorf("discovery.ScanPort: %w", err)
	}
	defer conn.Close()
	return true, nil
}

type Pair struct {
	Local  net.IP
	Remote net.IP
}

func Do(setting Setting) ([]Pair, error) {
	port := setting.Port.UnwrapOr(port.PORT + 1)
	limit := setting.Limit.UnwrapOr(3)
	limitTime := setting.LimitTime.UnwrapOr(time.Millisecond * 20)

	localIPs, err := ip.GetLocalIPs()
	if err != nil {
		return nil, err
	}

	addresses := make([]Pair, 0)

	for _, i := range localIPs {
		f, err := ip.GetFirstIP(i)
		if err != nil {
			continue
		}

		for {
			func() {
				if i.Equal(f) {
					return
				}
				b, err := ScanPort(f, port, limitTime)
				if err != nil {
					return
				}
				if b {
					addresses = append(addresses, Pair{
						Local:  i,
						Remote: f,
					})
				}
			}()
			if len(addresses) >= limit {
				return addresses, nil
			}
			n, err := ip.GetNextIP(f)
			if err != nil {
				break
			}
			f = n
		}
	}

	return addresses, nil
}

func Listen(setting Setting) error {
	port := setting.Port.UnwrapOr(port.PORT + 1)

	tcp, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return fmt.Errorf("discovery.Listen: %w", err)
	}

	for {
		conn, err := tcp.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		conn.Close()
	}
}
