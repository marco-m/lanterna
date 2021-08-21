package main

import (
	"fmt"
	"net"

	"github.com/rs/zerolog"
)

func cmdCollect(cfg Args) error {
	ips, err := collect(cfg.log)
	if err != nil {
		return err
	}
	for _, ip := range ips {
		fmt.Printf("%s\n", ip)
	}
	return nil
}

// collect returns a list of the global unicast IP addresses present on the host.
func collect(log zerolog.Logger) ([]string, error) {
	var ips []string

	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
			log.Err(err).Msg("iface.Addrs")
			continue
		}
		for _, addr := range addrs {
			ip, _, err := net.ParseCIDR(addr.String())
			if err != nil {
				log.Err(err).Stringer("addr", addr).Msg("ParseCIDR")
				continue
			}
			if !ip.IsGlobalUnicast() {
				continue
			}
			ips = append(ips, ip.String())
		}
	}

	if len(ips) == 0 {
		return nil, fmt.Errorf("could not find any IP address (%s)", ips)
	}

	return ips, nil
}
