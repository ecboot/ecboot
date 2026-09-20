package jobs

import "context"

type HelloJob struct {
}

func NewHelloJob() *HelloJob {
	return &HelloJob{}
}

func (s *HelloJob) Cron() string {
	var pattern = "@every 10m" // 每10分钟执行一次
	return pattern
}

func (s *HelloJob) Enable() bool {
	return true
}

func (s *HelloJob) Handle(ctx context.Context) {
}