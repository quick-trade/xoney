package trade

type Level interface {
	onUpdate()
	onBreakout()
	TradePart() float64
	TriggerPrice() float64
}
