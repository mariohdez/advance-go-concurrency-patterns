package timex

import "time"

//go:generate mockgen -package=$GOPACKAGE -source=$GOFILE -destination=mock_$GOFILE

type Ticker interface {
	C() <-chan time.Time
	Stop()
}

type TickerBuilder interface {
	NewTicker(d time.Duration) Ticker
}

type TickerFactory struct {
}

func (f *TickerFactory) NewTicker(d time.Duration) Ticker {
	return &TimeTicker{
		ticker: time.NewTicker(d),
	}
}

type TimeTicker struct {
	ticker *time.Ticker
}

func (t *TimeTicker) C() <-chan time.Time {
	return t.ticker.C
}

func (t *TimeTicker) Stop() {
	t.ticker.Stop()
}
