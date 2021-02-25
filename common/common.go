package common

import "time"

func Now() string {
	return time.Now().Format("2006-01-02 03:04:05")
}

type Timer struct {
	begin time.Time
}

func NewTimer() *Timer {
	return &Timer{time.Now()}
}

func (t *Timer) Restart() {
	t.begin = time.Now()
}

func (t *Timer) Elapse() time.Duration {
	return time.Since(t.begin)
}
