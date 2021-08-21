package main

import (
	"encoding/base32"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rs/xid"
	"github.com/shirou/gopsutil/v3/host"
)

func cmdRun(cfg Args) error {
	config, err := loadConfig(cfg.ConfigPath)
	if err != nil {
		return fmt.Errorf("run: %s", err)
	}

	sink := config.Sinks[0]

	ips, err := collect(cfg.log)
	if err != nil {
		return fmt.Errorf("run: %s", err)
	}

	var hostname string
	if hostname, err = os.Hostname(); err != nil {
		hostname = fmt.Sprintf("hostname: %s", err)
	}
	encoder := base32.StdEncoding.WithPadding(base32.NoPadding)
	machineID := encoder.EncodeToString(xid.New().Machine())

	bt, btErr := host.BootTime()
	bootTime := time.Unix(int64(bt), 0)
	now := time.Now().Truncate(time.Second)
	upDays := now.Sub(bootTime).Hours() / 24
	upTime := fmt.Sprintf("uptime: %.1f days (%s)", upDays, bootTime)
	if btErr != nil {
		upTime = fmt.Sprintf("uptime: %s", btErr)
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "*%s (%s)*\n", hostname, machineID)
	// fmt.Fprintf(&sb, "```\n")
	fmt.Fprintf(&sb, "now: %s\n", now)
	fmt.Fprintf(&sb, "%s\n", upTime)
	fmt.Fprintf(&sb, "IP addresses:\n")
	for _, ip := range ips {
		fmt.Fprintf(&sb, "    %s\n", ip)
	}
	// fmt.Fprintf(&sb, "```\n")

	// The gchat "threadKey" parameter will post all messages to the same thread.
	// https://developers.google.com/chat/reference/rest/v1/spaces.messages/create
	// The 3-byte (2^24=16_777_216) machine ID is unique enough for this usage.
	url := fmt.Sprintf("%s&threadKey=%s", sink.URL, machineID)

	cfg.log.Info().Msg("Sending message")
	return postJSON(url, map[string]string{"text": sb.String()})
}