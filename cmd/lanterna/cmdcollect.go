package main

import (
	"fmt"
	"log/slog"
	"net"
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
func collect(log *slog.Logger) ([]string, error) {
	var ips []string

	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
			log.Error("iface.Addrs", "err", err)
			continue
		}
		for _, addr := range addrs {
			ip, _, err := net.ParseCIDR(addr.String())
			if err != nil {
				log.Error("ParseCIDR", "addr", addr)
				continue
			}
			if !ip.IsGlobalUnicast() {
				continue
			}
			ips = append(ips, ip.String())
		}
	}
	return ips, nil
}
