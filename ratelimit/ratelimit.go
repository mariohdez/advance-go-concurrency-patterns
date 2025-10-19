package ratelimit

import "time"

type Service struct {
	bucket chan struct{}
	ticker *time.Ticker
	done   chan struct{}
}

func New(rate float64, limit int) *Service {
	srv := &Service{
		bucket: make(chan struct{}, limit),
		done:   make(chan struct{}),
		ticker: time.NewTicker(time.Duration(float64(time.Second) / rate)),
	}
	for range limit {
		srv.bucket <- struct{}{}
	}
	go func() {
		for {
			select {
			case <-srv.done:
				return
			case <-srv.ticker.C:
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
