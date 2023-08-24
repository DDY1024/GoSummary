package retry

import (
	"time"

	"github.com/DDY1024/GoSummary/library/rolling"
)

const (
	retryThreshold = 30 // 重试阈值
)

type Stats struct {
	success *rolling.Number // 成功次数
	fail    *rolling.Number // 失败次数(重试次数)
}

func newStats(name string) *Stats {
	return &Stats{
		success: rolling.NewNumber(),
		fail:    rolling.NewNumber(),
	}
}

func (m *Stats) shouldRetry() bool {
	now := time.Now()
	retryNum := m.fail.Sum(now)
	successNum := m.success.Sum(now)
	// 统计窗口内 rpc 重试次数超过阈值，或 重试次数 > 0.1 * 成功次数
	if retryNum > retryThreshold && retryNum > 0.1*successNum {
		return false
	}
	return true
}
