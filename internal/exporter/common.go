package exporter

import "github.com/go-sensors/rpi-sensor-exporter/internal/log"

var (
	shouldTerminate = func(err error) bool {
		if err == nil {
			return true
		}

		log.Warn("encountered recoverable error; allowing recovery to continue",
			"err", err)
		return false
	}
)
