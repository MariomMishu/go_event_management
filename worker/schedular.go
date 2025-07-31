package worker

import "time"

type Scheduler struct {
	ticker *time.Ticker // Ticker to send ticks at intervals
	done   chan bool    // A control channel used to gracefully stop the schedule
}

func NewScheduler(interval time.Duration) *Scheduler {
	return &Scheduler{
		ticker: time.NewTicker(interval),
		done:   make(chan bool),
	}
}

func (s *Scheduler) Start(task func()) {
	go func() {
		for {
			select {
			case <-s.done:
				return
			case <-s.ticker.C:
				task()
			}
		}
	}()
}
func (s *Scheduler) Stop() {
	s.ticker.Stop()
	s.done <- true
}
