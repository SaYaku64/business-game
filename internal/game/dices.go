package game

import (
	"math/rand"
	"time"
)

func genDices() (int, int) {
	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	first := r.Intn(6)
	second := r.Intn(6)

	return first, second
}

func sumDices() (sum int, isDouble bool) {
	first, second := genDices()

	sum = first + second
	isDouble = first == second

	return
}
