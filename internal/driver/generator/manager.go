package generator

import (
	"strings"
	"sync"
)

type GenManager struct {
	listGen map[string]*Generator
	mu      sync.Mutex
}

func CreateGenManager() *GenManager {
	return &GenManager{
		listGen: make(map[string]*Generator),
	}
}

func (g *GenManager) New(cfg string, flag bool) (*Generator, error) {
	if flag {
		if gen, ok := g.checkExistenceGen(cfg); ok {
			return gen, nil
		}
	}

	genType, genCfg, err := g.parseConfig(cfg)
	if err != nil {
		return nil, err
	}

	gen, err := g.selectGenType(genType, genCfg)
	if flag {
		g.mu.Lock()
		g.listGen[cfg] = gen
		g.mu.Unlock()
	}

	return gen, err
}

func (g *GenManager) parseConfig(cfg string) (string, string, error) {
	idx := strings.Index(cfg, delimiter)
	if idx < 0 {
		return "", "", errInvalidSettings
	}

	return cfg[:idx], cfg[idx+1:], nil
}

func (g *GenManager) selectGenType(genType, genCfg string) (*Generator, error) {
	switch genType {
	case sineGen:
		return NewSineGen(genCfg)
	case sawGen:
		return NewSawGen(genCfg)
	case randGen:
		return NewRandGen(genCfg)
	default:
		return nil, errGenTypeNotFound
	}
}

func (g *GenManager) checkExistenceGen(cfg string) (*Generator, bool) {
	g.mu.Lock()
	defer g.mu.Unlock()

	gen, ok := g.listGen[cfg]
	if ok {
		return gen, true
	}

	return nil, false
}
