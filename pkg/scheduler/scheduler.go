package scheduler

import (
	"time"
)

func Schedule(interval time.Duration, fn func()) {
	fn()
	ticker := time.NewTicker(interval)

	for {
		<-ticker.C
		fn()
	}
}
