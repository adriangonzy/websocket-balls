package game

import (
	"testing"
	"time"
)

func TestStartSimulation(t *testing.T) {
	Start()
	time.Sleep(20 * time.Second)
}
