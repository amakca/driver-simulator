package generator

import (
	m "practice/internal/models"
	"strconv"
	"strings"
	"time"
)

type sawSettings struct {
	amplitude float64
	frequency float64
}

func parseSawSettings(cfg string) (sawSettings, time.Duration, error) {
	parts := strings.Split(cfg, m.DELIMITER)
	if len(parts) != 3 {
		return sawSettings{}, 0, ErrInvalidSettings
	}

	settings := sawSettings{}
	sampleRate, _ := time.ParseDuration(parts[0])
	if sampleRate < MAX_SAMPLE_RATE {
		return sawSettings{}, 0, ErrSampleRateSmall
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
