package trade

type level interface {
	onUpdate()
	onBreakout()
	TradePart() float64
}

type callback func(*Level)

type Level struct {
	TriggerPrice float64
	tradePart    float64
	onUpdate     []callback
	onBreakout   []callback
}
