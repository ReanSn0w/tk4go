package tools

import (
	"time"
)

// Функция создает генератор кодов валидации
//
// ttl - время жизни кода в минутах
// len - длинна кода
// numbers, letters - указывает на то, какие символы можно использовать
// в генерируемых кодах
func NewCodeGenerator(ttl int, len int, numbers, letters bool) *CodeGenerator {
	return &CodeGenerator{
		randomizer: NewRandomizer(),
		numbers:    numbers,
		letters:    letters,
		ttl:        ttl + 2,
		len:        len,
	}
}

type CodeGenerator struct {
	randomizer *Randomizer

	numbers bool
	letters bool
	ttl     int
	len     int
}

func (cg *CodeGenerator) Generate(input string) string {
	return cg.GenerateForTime(time.Now(), input)
}

func (cg *CodeGenerator) GenerateForTime(t time.Time, input string) string {
	codes := cg.generateCodesArray(t.Add(time.Minute*time.Duration(cg.ttl-1)), input, 1)
	return codes[len(codes)-1]
}

func (cg *CodeGenerator) Validate(input, code string) bool {
	return cg.ValidateForTime(time.Now(), input, code)
}

func (cg *CodeGenerator) ValidateForTime(t time.Time, input string, code string) bool {
	codes := cg.generateCodesArray(t, input, cg.ttl)

	for _, c := range codes {
		if c == code {
			return true
		}
	}

	return false
}

func (cg *CodeGenerator) generateCodesArray(start time.Time, input string, limit int) []string {
	codes := make([]string, 0)
	inputInt := cg.inputToInt(input)

	for i := 0; i < limit; i++ {
		codeTime := start.Add(time.Minute * time.Duration(i))
		codeTime = time.Date(codeTime.Year(), codeTime.Month(), codeTime.Day(), codeTime.Hour(), codeTime.Minute(), 0, 0, codeTime.Location())
		codeTimeUnix := codeTime.Unix()
		code := cg.randomizer.PseudoString(codeTimeUnix+inputInt, cg.len, cg.letters, cg.numbers, false)
		codes = append(codes, code)
	}

	return codes
}

func (cg *CodeGenerator) inputToInt(input string) (result int64) {
	for i, val := range []rune(input) {
		result += int64(val) - int64(i)
	}

	return
}
