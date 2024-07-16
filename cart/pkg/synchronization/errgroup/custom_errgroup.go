package errgroup

import (
	"context"
	"sync"
	"time"
)

type ErrGroup struct {
	ctx       context.Context
	cancel    context.CancelFunc
	wg        sync.WaitGroup
	errCh     chan error
	semaphore *semaphore
}

// WithContext creates ErrGroup and returns  *ErrGroup and context.Context
func WithContext(ctx context.Context, rps int) (*ErrGroup, context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	return &ErrGroup{
		ctx:       ctx,
		cancel:    cancel,
		errCh:     make(chan error, 1),
		semaphore: makeSemaphore(rps),
	}, ctx
}

// Go runs f func in separate goroutine
func (g *ErrGroup) Go(f func() error) {
	g.wg.Add(1)

	go func() {
		defer g.wg.Done()

		select {
		case <-g.ctx.Done():
			return
		case <-g.semaphore.semCh: // if semaphore is available run f() function
			if err := f(); err != nil {
				select { // since errGroup should return only first err use select with default
				case g.errCh <- err: // add err to errCh and cancel context
					g.cancel()
				default:
				}
			}
		}
	}()
}

// Wait for all goroutines
func (g *ErrGroup) Wait() error {
	go func() {
		g.wg.Wait()        // wait for all goroutines
		close(g.errCh)     // close err chanel
		g.semaphore.Stop() // stop goroutine that controls semaphore
	}()

	select {
	case err := <-g.errCh:
		return err
	case <-g.ctx.Done():
		return g.ctx.Err()
	}
}

// semaphore is implemented using chanel. semCh chanel behaves like semaphore and rate-limiter at the same time
type semaphore struct {
	semCh  chan struct{}
	stopCh chan struct{}
	ticker *time.Ticker
}

func makeSemaphore(size int) *semaphore {
	stopCh := make(chan struct{}, 1)
	semCh := make(chan struct{}, size)
	ticker := time.NewTicker(time.Second)

	for i := 0; i < size; i++ {
		semCh <- struct{}{} // fill in semaphore
	}
	go func() {
		for {
			select {
			case <-stopCh:
				return
			case <-ticker.C: // every 1 second fill the channel with <rps> amount of data
				for i := 0; i < size; i++ {
					semCh <- struct{}{} // we will be locked if semaphore is full
				}
			}
		}
	}()
	return &semaphore{stopCh: stopCh, semCh: semCh, ticker: ticker}
}

func (s *semaphore) Stop() {
	close(s.stopCh)
	s.ticker.Stop()
}
