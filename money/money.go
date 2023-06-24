package money

import "math"

type Money int64

func (m Money) ToFloat() float64 {
	return float64(m) / 100 // Assuming the balance is represented in cents
}

func NewMoneyFromFloat(value float64) Money {
	amount := int64(math.Round(value * 100))
	return Money(amount)
}
