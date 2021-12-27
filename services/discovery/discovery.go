package discovery

import (
	"fmt"
	"log"
	"net"
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

func Do(setting *Setting) ([]string, error) {
	if setting == nil {
		setting = &Setting{
			Port:      port.PORT + 1,
			Limit:     3,
			LimitTime: time.Millisecond * 10,
		}
	}

	localIPs, err := ip.GetLocalIPs()
	if err != nil {
		return nil, err
	}

	addresses := make([]string, 0)

	for _, i := range localIPs {
		f, err := ip.GetFirstIP(i)
		if err != nil {
			log.Println(err)
			continue
		}

		for {
			func() {
				b, err := ScanPort(f, setting.Port, setting.LimitTime)
				if err != nil {
					log.Println(err)
					return
				}
				if b {
					addresses = append(addresses, f.String())
				}
			}()
			if len(addresses) >= setting.Limit {
				return addresses, nil
			}
			n, err := ip.GetNextIP(f)
			if err != nil {
				log.Println(err)
				break
			}
			f = n
		}
	}

	return addresses, nil
}
