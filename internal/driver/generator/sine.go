package generator

import (
	"math"
	"strconv"
	"strings"
	"time"
)

type sineSettings struct {
	amplitude float64
	frequency float64
}

func parseSineSettings(cfg string) (sineSettings, time.Duration, error) {
	parts := strings.Split(cfg, delimiter)
	if len(parts) != 3 {
		return sineSettings{}, 0, errInvalidSettings
	}

	settings := sineSettings{}
	sampleRate, _ := time.ParseDuration(parts[0])

	if sampleRate < maxPrescaler {
		return sineSettings{}, 0, errPrescallerSmall
	}
	settings.amplitude, _ = strconv.ParseFloat(parts[1], 64)
	settings.frequency, _ = strconv.ParseFloat(parts[2], 64)

	return settings, sampleRate, nil
}

// Конструктор синус-генератора
func NewSineGen(cfg string) (*Generator, error) {
	settings, sampleRate, err := parseSineSettings(cfg)
	if err != nil {
		return nil, err
	}

	return &Generator{
		valuer:     &settings,
		sampleRate: sampleRate,
	}, nil
}

func (s *sineSettings) value() float32 {
	return float32(s.amplitude * math.Sin(float64(
		time.Now().UnixNano())*2.0*math.Pi*
		float64(s.frequency)/1e9))
}
