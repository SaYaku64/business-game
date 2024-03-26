package game

import (
	"math/rand"
	"time"
)

func genDices() (int, int) {
	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	first := r.Intn(6) + 1
	second := r.Intn(6) + 1

	return first, second
}

func sumDices(first, second int) (sum int, isDouble bool) {
	return first + second, first == second
}
