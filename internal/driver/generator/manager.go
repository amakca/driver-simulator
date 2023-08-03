package generator

import (
	m "practice/internal/models"
	"strings"
	"sync"
)

type Manager struct {
	list map[string]*Generator
	rwmu sync.RWMutex
}

func CreateManager() (*Manager, error) {
	return &Manager{
		list: make(map[string]*Generator),
	}, nil
}

func (g *Manager) New(cfg string, useOptimization bool) (*Generator, error) {
	if useOptimization {
		if gen := g.findExistent(cfg); gen != nil {
			return gen, nil
		}
	}

	genType, genCfg, err := g.parseConfig(cfg)
	if err != nil {
		return nil, err
	}

	gen, err := g.selectGenType(genType, genCfg)
	if useOptimization {
		g.rwmu.Lock()
		g.list[cfg] = gen
		g.rwmu.Unlock()
	}

	return gen, err
}

func (g *Manager) parseConfig(cfg string) (string, string, error) {
	idx := strings.Index(cfg, m.DELIMITER)
	if idx < 0 {
		return "", "", ErrInvalidSettings
	}

	return strings.ToLower(cfg[:idx]), cfg[idx+1:], nil
}

func (g *Manager) selectGenType(genType, genCfg string) (*Generator, error) {
	switch genType {
	case SINE_GENERATOR:
		return NewSineGen(genCfg)
	case SAW_GENERATOR:
		return NewSawGen(genCfg)
	case RAND_GENERATOR:
		return NewRandGen(genCfg)
	default:
		return nil, ErrGenTypeNotFound
	}
}

func (g *Manager) findExistent(cfg string) *Generator {
	g.rwmu.RLock()
	defer g.rwmu.RUnlock()

	gen, ok := g.list[cfg]
	if ok {
		return gen
	}

	return nil
}
