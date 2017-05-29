package images

import (
	"math/rand"
	"sync"
	"time"
)

type lockedRandomSource struct {
	lk  sync.Mutex
	src rand.Source
}

var randomSource = rand.New(&lockedRandomSource{
	src: rand.NewSource(time.Now().UTC().UnixNano()),
})

func (r *lockedRandomSource) Int63() (n int64) {
	r.lk.Lock()
	n = r.src.Int63()
	r.lk.Unlock()
	return
}

func (r *lockedRandomSource) Seed(seed int64) {
	r.lk.Lock()
	r.src.Seed(seed)
	r.lk.Unlock()
}
