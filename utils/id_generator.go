package utils

import (
	"math/rand"
	"sync"
	"time"
)

type IDGenerator struct {
	mu   sync.Mutex
	rand *rand.Rand
}

func NewIDGenerator() *IDGenerator {
	return &IDGenerator{
		rand: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (g *IDGenerator) GenerateID() int {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.rand.Intn(900000000) + 100000000
}
