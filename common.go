package gsclient

import (
	"errors"
	"github.com/google/uuid"
	"time"
)

type isContinue func() (bool, error)

//isValidUUID validates the uuid
func isValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

//retryWithTimeout reruns a function within a period of time
func retryWithTimeout(targetFunc isContinue, timeout, delay time.Duration) error {
	timer := time.After(timeout)
	for {
		select {
		case <-timer:
			return errors.New("timeout reached when retrying")
		default:
			time.Sleep(delay) //delay between retries
			continueRetrying, err := targetFunc()
			if err != nil {
				return err
			}
			if !continueRetrying {
				return nil
			}
		}
	}
}

//retryWithLimitedNumOfRetries reruns a function within a number of retries
func retryWithLimitedNumOfRetries(targetFunc isContinue, numOfRetries int, delay time.Duration) error {
	retryNo := 0
	for retryNo <= numOfRetries {
		time.Sleep(delay) //delay between retries
		continueRetrying, err := targetFunc()
		if err != nil {
			return err
		}
		if !continueRetrying {
			return nil
		}
		retryNo++
	}
	return errors.New("maximum number of retries reached when retrying")
}
