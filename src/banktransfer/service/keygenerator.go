package service

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"sync/atomic"
)

// KeyGenerator generates keys in the pattern {hostname}-{count}
type KeyGenerator struct {
	counter int32
	hostName string
}

func NewKeyGenerator() *KeyGenerator {
	hostName, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}
	return &KeyGenerator{counter: 0, hostName: hostName}
}

func (g *KeyGenerator) getUniqueId() string {
	return fmt.Sprintf("%g-%d", g.hostName, atomic.AddInt32(&g.counter,1))
}


