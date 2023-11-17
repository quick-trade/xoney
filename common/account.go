package common

import "xoney/common/data"

type Portfolio struct {
	assets map[data.Currency]float64
}
func (p Portfolio) Total() float64 {
	panic("TODO: implement")
}
func (p Portfolio) Balance(currency data.Currency) float64 {
	return p.assets[currency]
}
