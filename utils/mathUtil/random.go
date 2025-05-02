package mathUtil

import (
	"math/rand"
	"time"
)

func Random(n int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(n)
}
