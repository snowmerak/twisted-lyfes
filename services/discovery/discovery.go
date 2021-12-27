package discovery

import (
	"net/http"

	"github.com/snowmerak/twisted-lyfes/net/ip"
	"github.com/snowmerak/twisted-lyfes/net/port"
)

type Setting struct {
	Port  string
	Limit int
}

func Do(setting *Setting) ([]string, error) {
	if setting == nil {
		setting = &Setting{
			Port:  port.PORT,
			Limit: 10,
		}
	}

	ips, err := ip.GetLocalIPs()
	if err != nil {
		return nil, err
	}

	var addresses []string
	for _, i := range ips {
		f, err := ip.GetFirstIP(i)
		if err != nil {
			return nil, err
		}
		i = f
		for {
			resp, err := http.Get("http://" + i.String() + ":" + setting.Port + "/discovery")
			if err != nil {
				break
			}
			defer resp.Body.Close()
			if resp.StatusCode == 200 {
				addresses = append(addresses, i.String())
				if len(addresses) >= setting.Limit {
					return addresses, nil
				}
			}
			n, err := ip.GetNextIP(i)
			if err != nil {
				break
			}
			i = n
		}
	}

	return addresses, nil
}
