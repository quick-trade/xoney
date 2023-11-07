package trade

import (
	"xoney/common/data"
)

type Trade struct {
	breakouts []level
	entries   []level
}

func (t *Trade) Update(candle data.Candle) {
	panic("TODO: implement")
}

func (t Trade) IsEqual(other *Trade) bool {
	panic("TODO: implement")
}
