package generator

import (
	"math/rand"
	m "practice/internal/models"
	"strconv"
	"strings"
	"time"
)

type random struct {
	high float64
	low  float64
}

func parseRandom(cfg string) (random, time.Duration, error) {
	parts := strings.Split(cfg, m.DELIMITER)
	if len(parts) != 3 {
		return random{}, 0, ErrInvalidSettings
	}

	valuer := random{}
	sampleRate, err := time.ParseDuration(parts[0])
	if err != nil {
		return random{}, 0, err
	}

	if sampleRate < MAX_SAMPLE_RATE {
		return random{}, 0, ErrSampleRateSmall
	}

	valuer.high, err = strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return random{}, 0, err
	}

	valuer.low, err = strconv.ParseFloat(parts[2], 64)
	if err != nil {
		return random{}, 0, err
	}

	return valuer, sampleRate, nil
}

// Конструктор рандом-генератора
func NewRandGen(cfg string) (*Generator, error) {
	valuer, sampleRate, err := parseRandom(cfg)
	if err != nil {
		return nil, err
	}

	return &Generator{
		valuer:     &valuer,
		sampleRate: sampleRate,
	}, nil
}

func (s *random) value() float32 {
	return float32((s.low) + rand.Float64()*(s.high-s.low))
}
