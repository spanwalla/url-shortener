package encoder

import (
	"math/rand"
	"sync"
)

type Random struct {
	Alphabet []rune
	Seed     int64
	rng      *rand.Rand
	mu       sync.Mutex
}

func NewRandom(alphabet string, seed int64) *Random {
	return &Random{
		Alphabet: []rune(alphabet),
		Seed:     seed,
		rng:      rand.New(rand.NewSource(seed)),
	}
}

func (r *Random) Encode(_ string, length int) string {
	if length <= 0 {
		return ""
	}

	buf := make([]rune, length)
	for i := range buf {
		r.mu.Lock()
		char := r.rng.Intn(len(r.Alphabet))
		r.mu.Unlock()
		buf[i] = r.Alphabet[char]
	}

	return string(buf)
}
