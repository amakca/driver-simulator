package generator

import (
	m "practice/internal/models"
	"strconv"
	"strings"
	"time"
)

type saw struct {
	amplitude float64
	frequency float64
}

func parseSaw(cfg string) (saw, time.Duration, error) {
	parts := strings.Split(cfg, m.DELIMITER)
	if len(parts) != 3 {
		return saw{}, 0, ErrInvalidSettings
	}

	valuer := saw{}
	sampleRate, err := time.ParseDuration(parts[0])
	if err != nil {
		return saw{}, 0, err
	}

	if sampleRate < MAX_SAMPLE_RATE {
		return saw{}, 0, ErrSampleRateSmall
	}
	valuer.amplitude, err = strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return saw{}, 0, err
	}

	valuer.frequency, err = strconv.ParseFloat(parts[2], 64)
	if err != nil {
		return saw{}, 0, err
	}

	return valuer, sampleRate, nil
}

// Конструктор пила-генератора
func NewSawGen(cfg string) (*Generator, error) {
	valuer, sampleRate, err := parseSaw(cfg)
	if err != nil {
		return nil, err
	}

	return &Generator{
		valuer:     &valuer,
		sampleRate: sampleRate,
	}, nil
}

func (s *saw) value() float32 {
	return float32(s.amplitude*(2*s.frequency*
		float64(time.Now().Second())) - s.amplitude)
}
