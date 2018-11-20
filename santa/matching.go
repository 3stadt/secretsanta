package santa

import (
	"math/rand"
	"time"
)

func Pair(p []participant) (*map[participant]participant, int64) {
	// shuffle the slice: https://en.wikipedia.org/wiki/Fisher%E2%80%93Yates_shuffle#The_modern_algorithm
	seed := time.Now().UnixNano()
	rand.Seed(seed)
	for i := len(p) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		p[i], p[j] = p[j], p[i]
	}
	// generate pairing
	lastIndex := len(p) - 1
	partMap := make(map[participant]participant)
	for i, part := range p {
		if i == lastIndex {
			partMap[part] = p[0]
			continue
		}
		partMap[part] = p[i+1]
	}
	return &partMap, seed
}
