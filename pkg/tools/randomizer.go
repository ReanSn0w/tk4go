package tools

import (
	"math/rand"
	"time"
)

const (
	alfabetValues = "QWERTYUIOPASDFGHJKLZXCVBNM"
	numbersValues = "1234567890"
	symbolsValues = "!@#$%^&*<>?"
)

func NewRandomizer() *Randomizer {
	return &Randomizer{}
}

type Randomizer struct{}

func (r *Randomizer) String(count int, alphabet, numbers, symbols bool) string {
	return r.PseudoString(time.Now().Unix(), count, alphabet, numbers, symbols)
}

func (r *Randomizer) PseudoString(seed int64, count int, alphabet, numbers, symbols bool) string {
	chars := ""

	if alphabet {
		chars += alfabetValues
	}

	if numbers {
		chars += numbersValues
	}

	if symbols {
		chars += symbolsValues
	}

	random := rand.New(rand.NewSource(seed))

	randomLimit := len(chars)
	values := make([]byte, count)
	index := 0

	for i := 0; i < count; i++ {
		index = random.Intn(randomLimit)
		values[i] = chars[index]
	}

	return string(values)
}
