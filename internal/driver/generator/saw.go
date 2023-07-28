package generator

import (
	"strconv"
	"strings"
	"time"
)

type sawSettings struct {
	amplitude float64
	frequency float64
}

func parseSawSettings(cfg string) (sawSettings, time.Duration, error) {
	parts := strings.Split(cfg, delimiter)
	if len(parts) != 3 {
		return sawSettings{}, 0, errInvalidSettings
	}

	settings := sawSettings{}
	sampleRate, _ := time.ParseDuration(parts[0])
	if sampleRate < maxPrescaler {
		return sawSettings{}, 0, errPrescallerSmall
	}
	settings.amplitude, _ = strconv.ParseFloat(parts[1], 64)
	settings.frequency, _ = strconv.ParseFloat(parts[2], 64)

	return settings, sampleRate, nil
}

// Конструктор пила-генератора
func NewSawGen(cfg string) (*Generator, error) {
	settings, sampleRate, err := parseSawSettings(cfg)
	if err != nil {
		return nil, err
	}

	return &Generator{
		valuer:     &settings,
		sampleRate: sampleRate,
	}, nil
}

func (s *sawSettings) value() float32 {
	return float32(s.amplitude*(2*s.frequency*
		float64(time.Now().Second())) - s.amplitude)
}
