package generator

import (
	"bytes"
	"encoding/binary"
	"math"
	m "practice/internal/models"
	u "practice/internal/utils"
	"time"
)

type Generator struct {
	valuer     valuer
	subs       uint32
	value      float32
	sampleRate time.Duration
	done       chan struct{}
}

type valuer interface {
	value() float32
}

func (g *Generator) Start() error {
	if g.subs == 0 {
		g.done = make(chan struct{})
		g.subs = 1
		go func() {
			ticker := time.NewTicker(g.sampleRate)
			for {
				select {
				case <-ticker.C:
					g.value = g.valuer.value()

				case <-g.done:
					ticker.Stop()
					return
				}
			}
		}()
	} else {
		g.subs++
	}
	return nil
}

func (g *Generator) Stop() error {
	switch g.subs {
	case 0:
		return ErrGenAlreadyStopped
	case 1:
		if !u.IsChanClosable(g.done) {
			return m.ErrCannotCloseChan
		}
		close(g.done)
		g.subs = 0
		return nil
	default:
		g.subs--
		return nil
	}
}

func (g *Generator) ValueBytes() []byte {
	bits := math.Float32bits(g.value)
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, bits)

	return bytes
}

func (g *Generator) Value() float32 {
	return g.value
}

func (g *Generator) SetValue(value float32) {
	g.value = value
}

func (g *Generator) SetValueBytes(value []byte) error {
	buf := bytes.NewReader(value)
	err := binary.Read(buf, binary.LittleEndian, &g.value)
	return err
}
