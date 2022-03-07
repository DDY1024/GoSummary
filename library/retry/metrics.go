package retry

import (
	"sync"
	"time"

	"github.com/DDY1024/GoSummary/library/rolling"
)

const (
	retryThreshold = 30 // 重试阈值
)

type metrics struct {
	success *rolling.Number // 成功次数
	fail    *rolling.Number // 失败次数(重试次数)

	// 暂时不用，用于
	// tagkv map[string]string  //
}

func newMetrics(name string) *metrics {
	return &metrics{
		success: rolling.NewNumber(),
		fail:    rolling.NewNumber(),
	}
}

func (m *metrics) shouldRetry() bool {
	now := time.Now()
	retryNum := m.fail.Sum(now)
	successNum := m.success.Sum(now)
	// 统计窗口内 rpc 重试次数超过阈值，或 重试次数 > 0.1 * 成功次数
	if retryNum > retryThreshold && retryNum > 0.1*successNum {
		return false
	}
	return true
}

var retryMetrics = make(map[string]*metrics)
var retryMetricsLock sync.RWMutex

func getMetrics(name string) *metrics {
	retryMetricsLock.RLock()
	if m, ok := retryMetrics[name]; ok {
		retryMetricsLock.RUnlock()
		return m
	}
	retryMetricsLock.RUnlock()

	retryMetricsLock.Lock()
	defer retryMetricsLock.Unlock()

	if m, ok := retryMetrics[name]; ok {
		return m
	}

	retryMetrics[name] = newMetrics(name)
	return retryMetrics[name]
}
