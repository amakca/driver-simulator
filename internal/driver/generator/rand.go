package generator

import (
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type randSettings struct {
	high float64
	low  float64
}

func parseRandSettings(cfg string) (randSettings, time.Duration, error) {
	parts := strings.Split(cfg, delimiter)
	if len(parts) != 3 {
		return randSettings{}, 0, errInvalidSettings
	}

	settings := randSettings{}
	sampleRate, _ := time.ParseDuration(parts[0])
	if sampleRate < maxPrescaler {
		return randSettings{}, 0, errPrescallerSmall
	}
	settings.high, _ = strconv.ParseFloat(parts[1], 64)
	settings.low, _ = strconv.ParseFloat(parts[2], 64)

	return settings, sampleRate, nil
}

// Конструктор рандом-генератора
func NewRandGen(cfg string) (*Generator, error) {
	settings, sampleRate, err := parseRandSettings(cfg)
	if err != nil {
		return nil, err
	}

	return &Generator{
		valuer:     &settings,
		sampleRate: sampleRate,
	}, nil
}

func (s *randSettings) value() float32 {
	return float32((s.low) + rand.Float64()*(s.high-s.low))
}
