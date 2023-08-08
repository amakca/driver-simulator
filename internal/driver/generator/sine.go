package generator

import (
	"math"
	m "practice/internal/models"
	"strconv"
	"strings"
	"time"
)

type sine struct {
	amplitude float64
	frequency float64
}

func parseSine(cfg string) (sine, time.Duration, error) {
	parts := strings.Split(cfg, m.DELIMITER)
	if len(parts) != 3 {
		return sine{}, 0, ErrInvalidSettings
	}

	valuer := sine{}
	sampleRate, err := time.ParseDuration(parts[0])
	if err != nil {
		return sine{}, 0, err
	}

	if sampleRate < MAX_SAMPLE_RATE {
		return sine{}, 0, ErrSampleRateSmall
	}

	valuer.amplitude, err = strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return sine{}, 0, err
	}

	valuer.frequency, err = strconv.ParseFloat(parts[2], 64)
	if err != nil {
		return sine{}, 0, err
	}

	return valuer, sampleRate, nil
}

// Конструктор синус-генератора
func NewSineGen(cfg string) (*Generator, error) {
	valuer, sampleRate, err := parseSine(cfg)
	if err != nil {
		return nil, err
	}

	return &Generator{
		valuer:     &valuer,
		sampleRate: sampleRate,
	}, nil
}

func (s *sine) value() float32 {
	return float32(s.amplitude * math.Sin(2.0*math.Pi*s.frequency*
		float64(time.Now().Second())))
}
