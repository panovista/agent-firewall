package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

var (
	currentStreams int64
	peakStreams    int64
	streamMu       sync.RWMutex
	bootTime       time.Time
	nodeID         string
	activeLicense  string
)

type PanovistaMetric struct {
	Level               string `json:"level"`
	Tag                 string `json:"tag"`
	Status              string `json:"status"`
	NodeID              string `json:"node_id"`
	LicenseTierClaimed  string `json:"license_tier_claimed"`
	UptimeSeconds       int64  `json:"uptime_seconds"`
	PeakConcurrent      int64  `json:"peak_concurrent_streams"`
}

func initTelemetry(licenseTier string) {
	bootTime = time.Now()
	nodeID = os.Getenv("HOSTNAME")
	if nodeID == "" {
		nodeID = "unknown-node"
	}
	activeLicense = licenseTier
}

func trackStreamStart() {
	streamMu.Lock()
	defer streamMu.Unlock()
	currentStreams++
	if currentStreams > peakStreams {
		peakStreams = currentStreams
	}
}

func trackStreamEnd() {
	streamMu.Lock()
	defer streamMu.Unlock()
	currentStreams--
}

// emitMetric executes the Phase 1, 2, and 3 standard output stamping
func emitMetric(status string) {
	streamMu.RLock()
	peak := peakStreams
	streamMu.RUnlock()

	uptime := int64(time.Since(bootTime).Seconds())
	if status == "boot" {
		uptime = 0
	}

	metric := PanovistaMetric{
		Level:               "info",
		Tag:                 "PANOVISTA_METRIC",
		Status:              status,
		NodeID:              nodeID,
		LicenseTierClaimed:  activeLicense,
		UptimeSeconds:       uptime,
		PeakConcurrent:      peak,
	}

	out, _ := json.Marshal(metric)
	fmt.Fprintln(os.Stdout, string(out))
}

// Phase 2: Automated 24-Hour Interval
func startTelemetryTicker() {
	ticker := time.NewTicker(24 * time.Hour)
	go func() {
		for range ticker.C {
			emitMetric("active")
		}
	}()
}