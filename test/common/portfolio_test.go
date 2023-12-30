package common_test

import (
	"testing"

	"xoney/common"
	"xoney/common/data"
)

func usd() data.Currency {
	return data.NewCurrency("USD", "NYSE")
}
func portfolio() common.Portfolio {
	return common.NewPortfolio(usd())
}

func TestSet(t *testing.T) {
	USD :=usd()
	p := portfolio()
	p.Set(USD, 100)

	if p.Balance(USD) != 100 {
		t.Error("Portfolio.Set() should impact portfolio.Balance")
	}
}

func TestCopy(t *testing.T) {
	USD :=usd()
	p := portfolio()
	p.Set(USD, 100)

	pc := p.Copy()
	p.Set(USD, 50)

	if pc.Balance(USD) != 100 {
		t.Error("Portfolio.Copy() is not working properly")
	}
}
