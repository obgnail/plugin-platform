package math

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

func MatchWildcard(input string, wildcard string) bool {
	runeInput := []rune(input)
	runePattern := []rune(wildcard)
	lenInput := len(runeInput)
	lenPattern := len(runePattern)

	isMatchingMatrix := make([][]bool, lenInput+1)

	for i := range isMatchingMatrix {
		isMatchingMatrix[i] = make([]bool, lenPattern+1)
	}

	isMatchingMatrix[0][0] = true
	for i := 1; i < lenInput; i++ {
		isMatchingMatrix[i][0] = false
	}

	if lenPattern > 0 {
		if runePattern[0] == '*' {
			isMatchingMatrix[0][1] = true
		}
	}

	for j := 2; j <= lenPattern; j++ {
		if runePattern[j-1] == '*' {
			isMatchingMatrix[0][j] = isMatchingMatrix[0][j-1]
		}

	}

	for i := 1; i <= lenInput; i++ {
		for j := 1; j <= lenPattern; j++ {

			if runePattern[j-1] == '*' {
				isMatchingMatrix[i][j] = isMatchingMatrix[i-1][j] || isMatchingMatrix[i][j-1]
			}

			if runePattern[j-1] == '?' || runeInput[i-1] == runePattern[j-1] {
				isMatchingMatrix[i][j] = isMatchingMatrix[i-1][j-1]
			}
		}
	}

	return isMatchingMatrix[lenInput][lenPattern]
}
