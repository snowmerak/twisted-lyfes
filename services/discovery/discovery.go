package discovery

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/snowmerak/twisted-lyfes/net/ip"
	"github.com/snowmerak/twisted-lyfes/net/port"
)

type Setting struct {
	Port      int
	Limit     int
	LimitTime time.Duration
}

func ScanPort(ip net.IP, port int, limit time.Duration) (bool, error) {
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

func Do(setting *Setting) ([]Pair, error) {
	if setting == nil {
		setting = &Setting{
			Port:      port.PORT + 1,
			Limit:     3,
			LimitTime: time.Millisecond * 20,
		}
	}

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
				b, err := ScanPort(f, setting.Port, setting.LimitTime)
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
			if len(addresses) >= setting.Limit {
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

func Listen(setting *Setting) error {
	if setting == nil {
		setting = &Setting{
			Port: port.PORT + 1,
		}
	}

	tcp, err := net.Listen("tcp", ":"+strconv.Itoa(setting.Port))
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
