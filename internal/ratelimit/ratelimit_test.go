package ratelimit_test

import (
	"advance-go-concurrency-patterns/internal/ratelimit"
	"advance-go-concurrency-patterns/internal/timex"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type rateLimitSuite struct {
	suite.Suite

	ctrl *gomock.Controller

	ticker      *timex.MockTicker
	tickCh      chan time.Time
	rateLimiter *ratelimit.Service
}

func (s *rateLimitSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())

	s.ticker = timex.NewMockTicker(s.ctrl)
	s.tickCh = make(chan time.Time)
	s.ticker.EXPECT().C().Return(s.tickCh).AnyTimes()
	s.rateLimiter = ratelimit.New(5, s.ticker)
}

func (s *rateLimitSuite) SetupSubTest() {
	s.SetupTest()
}

func (s *rateLimitSuite) TearDownTest() {
	s.ticker.EXPECT().Stop()
	s.rateLimiter.Close()

	s.ctrl.Finish()
}

func (s *rateLimitSuite) TearDownSubTest() {
	s.TearDownTest()
}

func (s *rateLimitSuite) TestWait() {
	s.Run("when no rehydration of bucket", func() {
		s.rateLimiter.Wait()
		s.rateLimiter.Wait()
		s.rateLimiter.Wait()
		s.rateLimiter.Wait()
		s.rateLimiter.Wait()
	})

	s.Run("when need to rehydrate tokens", func() {
		s.rateLimiter.Wait()
		s.rateLimiter.Wait()
		s.rateLimiter.Wait()
		s.rateLimiter.Wait()
		s.rateLimiter.Wait()

		var wg sync.WaitGroup
		blocked := make(chan struct{})
		wg.Add(1)
		go func() {
			defer wg.Done()

			close(blocked)
			s.rateLimiter.Wait()
		}()

		<-blocked
		s.tickCh <- time.Time{}

		done := make(chan struct{})
		go func() {
			wg.Wait()

			close(done)
		}()

		select {
		case <-done:
			// success
		case <-time.After(5 * time.Second):
			s.Require().FailNow("deadlock?")
		}
	})
}

func TestRateLimitSuite(t *testing.T) {
	suite.Run(t, &rateLimitSuite{})
}
