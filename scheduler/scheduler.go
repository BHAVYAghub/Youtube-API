package scheduler

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	log "github.com/BHAVYAghub/Youtube-API/logging"
	"go.uber.org/zap"
)

// SubmitToTicker submits a func to be called to a ticker
func SubmitToTicker(wg *sync.WaitGroup, tickerFunc func() error, period time.Duration) {
	ticker := time.NewTicker(period)

	wg.Add(1)
	go initTicker(wg, ticker, tickerFunc)
	wg.Wait()
}

func initTicker(wg *sync.WaitGroup, ticker *time.Ticker, tickerFunc func() error) {
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	for {
		select {
		case <-ticker.C:
			currentTime := time.Now()
			log.Info("Calling ticker function")

			if err := tickerFunc(); err != nil {
				log.Error("Scheduled task via ticker failed", zap.Any("time_taken", time.Since(currentTime)), zap.Error(err))
			}
		case <-quit:
			ticker.Stop()
			wg.Done()
			return
		}
	}
}
