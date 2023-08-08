package main

import (
	"encoding/base32"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/rs/xid"
	"github.com/shirou/gopsutil/v3/host"
)

func cmdRun(args Args) error {
	config, err := loadConfig(args.ConfigPath)
	if err != nil {
		return fmt.Errorf("run: %s", err)
	}
	encoder := base32.StdEncoding.WithPadding(base32.NoPadding)
	machineID := encoder.EncodeToString(xid.New().Machine())

	ticker := time.NewTicker(24 * time.Hour)

	for {
		if err := runHandle(args, config, machineID, collect, postJSON); err != nil {
			return err
		}
		t := <-ticker.C
		args.log.Info("ticker", "tick", t)
	}
}

type collectFn func(log *slog.Logger) ([]string, error)

type postJSONFn func(url string, msg map[string]string) error

func runHandle(args Args, config configuration, machineID string, collect collectFn, postJSON postJSONFn) error {
	sink := config.Sinks[0]

	ips, err := collect(args.log)
	if err != nil {
		return fmt.Errorf("run: %s", err)
	}
	if len(ips) == 0 {
		args.log.Warn("could not find any IP address")
	}

	var hostname string
	if hostname, err = os.Hostname(); err != nil {
		hostname = fmt.Sprintf("hostname: %s", err)
		args.log.Error("could not find hostname")
	}

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
	if len(ips) == 0 {
		fmt.Fprintf(&sb, "    WARNING: none found\n")
	}
	// fmt.Fprintf(&sb, "```\n")

	// The gchat "threadKey" parameter will post all messages to the same thread.
	// https://developers.google.com/chat/reference/rest/v1/spaces.messages/create
	// The 3-byte (2^24=16_777_216) machine ID is unique enough for this usage.
	url := fmt.Sprintf("%s&threadKey=%s", sink.URL, machineID)

	args.log.Info("Sending message")
	err = postJSON(url, map[string]string{"text": sb.String()})
	if err != nil {
		args.log.Error("postJSON", "err", err)
	}
	return nil
}
