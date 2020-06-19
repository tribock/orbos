package helpers

import (
	"time"

	"github.com/pkg/errors"
)

func Retry(timer *time.Timer, interval time.Duration, callback func() (retry bool)) error {
	if !callback() {
		if !timer.Stop() {
			<-timer.C
		}
		return nil
	}

	for {
		select {
		case <-timer.C:
			return errors.New("Timed out")
		case <-time.After(interval):
			if !callback() {
				if !timer.Stop() {
					<-timer.C
				}
				return nil
			}
			continue
		}
	}
}
