package retry

import (
	"errors"
	"sync/atomic"
	"time"
)

const (
	retryMetricsName        = "retry.throughput"
	retrySuccessMetricsName = "retry.success.throughput"
)

var (
	Retry0TimeError    = errors.New("retry time can't be 0")
	RetryTimeout0Error = errors.New("retry timeout can't be 0")
	TimeoutErr         = errors.New("retry timeout")
	ShouldNotRetryErr  = errors.New("Shouldn't retry")
)

type runFunc func() error

func Do(name string, times int, timeout time.Duration, run runFunc) error {
	// 重试次数
	if times == 0 {
		return Retry0TimeError
	}

	// 超时时间
	if timeout == 0 {
		return RetryTimeout0Error
	}

	var (
		err          error
		timeoutAbort int64
		counter      = getMetrics(name)
		resultCh     = make(chan error, 1)
		panicCh      = make(chan struct{}, 1)
	)

	go func() {
		var runErr error
		defer func() {
			if err := recover(); err != nil {
				// FIXME: log
				// fmt.Println("panic err, stack=", string(debug.Stack()))
				panicCh <- struct{}{}
			}
		}()
		for i := 0; i < times; i++ {
			// 超过总超时设置，直接返回，避免无效调用
			if atomic.LoadInt64(&timeoutAbort) > 0 {
				return
			}

			if i > 0 && !counter.shouldRetry() {
				resultCh <- ShouldNotRetryErr
				return
			}

			runErr = run()
			if runErr == nil {
				resultCh <- nil
				counter.success.Increment(1)
				return
			}
			counter.fail.Increment(1)
		}
		resultCh <- runErr
		return
	}()

	timer := time.NewTimer(timeout)
	defer timer.Stop()
	select {
	case err = <-resultCh:
		// do err log
	case <-panicCh:
		// do panic process
	case <-timer.C:
		// do timeout process
		err = TimeoutErr
		atomic.StoreInt64(&timeoutAbort, 1)
	}
	return err
}
