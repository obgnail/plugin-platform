package seq

import (
	crand "crypto/rand"
	"encoding/binary"
	"github.com/obgnail/plugin-platform/common/log"
	"math/rand"
)

func CreateCaptcha() uint64 {
	var src cryptoSource
	rnd := rand.New(src)
	re := rnd.Intn(100000)
	return uint64(re)
}

type cryptoSource struct{}

func (s cryptoSource) Seed(seed int64) {}

func (s cryptoSource) Int63() int64 {
	return int64(s.Uint64() & ^uint64(1<<63))
}

func (s cryptoSource) Uint64() (v uint64) {
	err := binary.Read(crand.Reader, binary.BigEndian, &v)
	if err != nil {
		log.ErrorDetails(err)
	}
	return v
}
