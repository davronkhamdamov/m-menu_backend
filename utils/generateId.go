package utils

import (
	"fmt"
	"math/rand"
	"time"
)

func GenerateOrderID() string {
	rand.Seed(time.Now().UnixNano())
	first := rand.Intn(10000)
	second := rand.Intn(10000)
	return fmt.Sprintf("%04d-%04d", first, second)
}
