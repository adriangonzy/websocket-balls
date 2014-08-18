package game

import (
	"fmt"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

func randFloat(min, max int) float64 {
	return rand.Float64()*float64(max-min) + float64(min)
}

func randInt(min, max int) int {
	return rand.Intn(max-min) + min
}

func randomColor() string {
	return fmt.Sprintf("#%x", uint(rand.Float64()*float64(0xffffff)))
}
