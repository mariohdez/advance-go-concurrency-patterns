package ratelimit

import "time"

type Ticker interface {
	C() <-chan time.Time
	Stop()
}

type Service struct {
	bucket chan struct{}
	done   chan struct{}
	ticker Ticker
}

func New(rate float64, limit int, ticker Ticker) *Service {
	srv := &Service{
		bucket: make(chan struct{}, limit),
		done:   make(chan struct{}),
		ticker: ticker,
	}
	for range limit {
		srv.bucket <- struct{}{}
	}
	go func() {
		for {
			select {
			case <-srv.done:
				return
			case <-srv.ticker.C():
				select {
				case srv.bucket <- struct{}{}:
				default:
				}
			}
		}
	}()
	return srv
}

func (s *Service) Wait() {
	<-s.bucket
}

func (s *Service) Close() {
	s.done <- struct{}{}
	s.ticker.Stop()
}
