package generator

import (
	"math/rand"
	"time"
)

const (
	AlphaLetters  = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	NumericDigits = "0123456789"
	AlphaNumeric  = AlphaLetters + NumericDigits
)

type CharsetGenerator struct {
	rnd *rand.Rand
}

func NewCharsetGenerator() *CharsetGenerator {
	return &CharsetGenerator{
		rnd: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Generate возвращает случайную строку длины n из букв и цифр.
func (g *CharsetGenerator) Generate(n int) string {
	return g.GenerateCustom(n, AlphaNumeric)
}

// GenerateAlpha возвращает строку длины n, содержащую только буквы.
func (g *CharsetGenerator) GenerateAlpha(n int) string {
	return g.GenerateCustom(n, AlphaLetters)
}

// GenerateNumeric возвращает строку длины n, содержащую только цифры.
func (g *CharsetGenerator) GenerateNumeric(n int) string {
	return g.GenerateCustom(n, NumericDigits)
}

// GenerateAlphaNum то же, что Generate (буквы+цифры).
func (g *CharsetGenerator) GenerateAlphaNum(n int) string {
	return g.Generate(n)
}

// GenerateCustom возвращает строку длины n из символов заданного набора.
func (g *CharsetGenerator) GenerateCustom(n int, charset string) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = charset[g.rnd.Intn(len(charset))]
	}
	return string(b)
}
