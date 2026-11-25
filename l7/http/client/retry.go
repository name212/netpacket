package client

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"
)

const (
	RetryDefaultSleep = 5 * time.Second
)

var (
	ErrRetryStop             = fmt.Errorf("stop retry")
	ErrRetryAttemptsExceeded = fmt.Errorf("all attempts exceeded")
)

type RetryFunc func(err error, attempt int, curSleep time.Duration, statusCode int, logger *slog.Logger) (time.Duration, error)

type Retrier struct {
	f      RetryFunc
	logger *slog.Logger

	timer    *time.Timer
	curSleep time.Duration
	attempt  int
}

func NewRetrier(f RetryFunc) *Retrier {
	if f == nil {
		return nil
	}

	return &Retrier{
		f: f,
	}
}

func (r *Retrier) WithLogger(logger *slog.Logger) *Retrier {
	if r != nil && logger != nil {
		r.logger = logger
	}

	return r
}

func (r *Retrier) makeRetry(ctx context.Context, err error, statusCode int) error {
	if r == nil {
		return ErrRetryStop
	}

	r.attempt += 1

	var errR error

	r.curSleep, errR = r.f(err, r.attempt, r.curSleep, statusCode, r.logger)

	if errR != nil {
		resErr := errR
		if !errors.Is(errR, ErrRetryStop) {
			resErr = fmt.Errorf("%w: %w", ErrRetryStop, errR)
		}

		r.stopTimer()

		return resErr
	}

	if r.curSleep > 0 {
		select {
		case <-ctx.Done():
			r.stopTimer()
			return fmt.Errorf("%w: context canceled with '%w' when retry on %w", ErrRetryStop, ctx.Err(), err)
		case <-r.getTimer().C:
			return nil
		}
	}

	return nil
}

func (r *Retrier) getTimer() *time.Timer {
	if r.timer != nil {
		r.timer.Reset(r.curSleep)
		return r.timer
	} else {
		r.timer = time.NewTimer(r.curSleep)
	}

	return r.timer
}

func (r *Retrier) stopTimer() {
	if r.timer != nil {
		r.timer.Stop()
		r.timer = nil
	}
}

func NewAttemptsRetries(attempts int, sleep time.Duration, logConsumer LogFuncConsumer, retryStatusCodes ...int) *Retrier {
	if attempts < 1 {
		return nil
	}

	if sleep == 0 {
		sleep = DefaultTimeoutSeconds
	}

	sleepStr := sleep.String()

	codes := make(map[int]struct{}, len(retryStatusCodes))
	for _, c := range retryStatusCodes {
		codes[c] = struct{}{}
	}

	var lastErr error

	setLastErr := func(err error, code int) {
		res := err
		if err == nil {
			res = fmt.Errorf("%w: %d", ErrInvalidStatusCode, code)
		}
		lastErr = res
	}

	f := func(err error, attempt int, _ time.Duration, statusCode int, logger *slog.Logger) (time.Duration, error) {
		setLastErr(err, statusCode)
		if attempt > attempts {
			return 0, fmt.Errorf("%w: %w %d last error %w", ErrRetryStop, ErrRetryAttemptsExceeded, attempts, lastErr)
		}

		if err == nil && statusCode > 0 {
			_, ok := codes[statusCode]
			if !ok {
				return 0, fmt.Errorf("%w: last error %w", ErrRetryStop, lastErr)
			}
		}

		if logger != nil {
			lf := logConsumer(logger)
			lastErrString := lastErr.Error()
			lf(
				"request failed do retry with sleep",
				slog.Int(LogRetryAttemptAttr, attempt),
				slog.String(LogRetrySleepAttr, sleepStr),
				slog.String(LogRetryErrAttr, lastErrString),
				slog.String(LogErrAttr, lastErrString),
				slog.Int(LogStatusCodeAttr, statusCode),
			)
		}

		return sleep, nil
	}

	return NewRetrier(f)
}
