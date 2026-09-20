package scheduler

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync/atomic"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcron"
	"github.com/redis/go-redis/v9"
)

// 集群调度器：多实例部署时保证 Scanner/Worker 只在一个实例上运行。
//
// 选主模型：实例通过 SET NX EX 抢占 Redis 领导锁（值为本实例令牌），
// 持锁期间周期性用 CAS 脚本续约（GET==token 才 EXPIRE，防误续他人的锁）；
// 续约发现锁易主（含锁过期被抢占后 Redis 恢复）即摘除本机全部任务，
// 保证任意时刻最多一个实例注册后台任务。主实例宕机后锁到期（≤TTL），
// 其余实例在下个检查周期接管。
//
// 任务摘除用 gcron.Remove 按名移除，重新当选时再注册——
// 任务本身（Stream 消费组、守卫更新、订单锁）本已幂等，选主解决的是
// 多实例重复扫描/空转的资源浪费与告警噪音，而非正确性问题。

const (
	// schedulerLockKey 领导锁 key。
	schedulerLockKey = "lock:scheduler:leader"
	// schedulerLockTTLSeconds 锁 TTL：宕机接管的最坏延迟 ≈ TTL + 检查周期。
	schedulerLockTTLSeconds = 30
)

// schedulerCheckInterval 选主检查周期（抢锁/续约共用），测试可调小。
var schedulerCheckInterval = 10 * time.Second

// schedulerJob 一个受选主管控的后台任务：name 与 gcron 注册名一致（摘除用）。
type schedulerJob struct {
	name  string
	start func(ctx context.Context)
}

// schedulerJobs 全部 Scanner/Worker 后台任务，仅选主实例注册。
// 声明为变量供测试替换为假任务（与 createPayment/newWechatClient 同模式）。
var schedulerJobs = []schedulerJob{
	// {"invoice_reissue_worker", invoice.StartReissueWorker},
}

// renewLockScript 续约 CAS：仅当锁值仍是本实例令牌时续期，否则返回 0（锁已易主）。
var renewLockScript = redis.NewScript(`
if redis.call("GET", KEYS[1]) == ARGV[1] then
	return redis.call("EXPIRE", KEYS[1], ARGV[2])
else
	return 0
end
`)

// releaseLockScript 释放 CAS：仅删除本实例持有的锁。
var releaseLockScript = redis.NewScript(`
if redis.call("GET", KEYS[1]) == ARGV[1] then
	return redis.call("DEL", KEYS[1])
else
	return 0
end
`)

type schedulerProvider struct {
	leader   atomic.Bool
	resigned atomic.Bool // 让位后终止：不再抢锁、选主循环退出（进程级终态）
	token    string
	stopped  chan struct{} // electionLoop 退出时关闭（测试据此等待 goroutine 终止）
}

var Scheduler = &schedulerProvider{}

// Boot 启动选主循环（后台 goroutine，不阻塞启动）。
func (a *schedulerProvider) Boot(ctx context.Context) {
	host, _ := os.Hostname()
	a.token = fmt.Sprintf("%s-%d-%d", host, os.Getpid(), time.Now().UnixNano())
	a.stopped = make(chan struct{})
	go a.electionLoop(ctx)
}

// IsLeader 当前实例是否为调度主实例（测试与诊断用）。
func (a *schedulerProvider) IsLeader() bool {
	return a.leader.Load()
}

func (a *schedulerProvider) electionLoop(ctx context.Context) {
	defer close(a.stopped)
	a.tick(ctx)
	ticker := time.NewTicker(schedulerCheckInterval)
	defer ticker.Stop()
	for {
		if a.resigned.Load() {
			// 显式让位（resign）后退出循环——不允许刚释放锁的实例立刻又选自己。
			return
		}
		select {
		case <-ctx.Done():
			a.resign(ctx)
			return
		case <-ticker.C:
			a.tick(ctx)
		}
	}
}

// universalClient 获取通用 Redis 客户端。
func (a *schedulerProvider) universalClient() (redis.UniversalClient, bool) {
	gredis := g.Redis()
	// 获取原始客户端
    client := gredis.Client()
    // 类型断言为 UniversalClient
    universalClient, ok := client.(redis.UniversalClient)
    if !ok {
        return nil, false
    }
    return universalClient, true
}

// tick 执行一次选主检查：非主抢锁，主续约；状态翻转时注册/摘除任务。
func (a *schedulerProvider) tick(ctx context.Context) {
	if a.resigned.Load() {
		return
	}
	if a.IsLeader() {
		held, err := a.renew(ctx)
		switch {
		case err != nil:
			// Redis 抖动时保持现状下个周期再试，不抖动任务注册；
			// 若锁真的易主，恢复后的续约 CAS 会返回未持有并在此摘除。
			g.Log().Warning(ctx, "调度器：续约领导者锁失败，保持当前状态并重试", "err", err)
		case !held:
			a.loseLeadership(ctx)
		}
		return
	}
	acquired, err := a.acquire(ctx)
	if err != nil {
		g.Log().Warning(ctx, "调度器：获取领导者锁失败，将在下一个周期重试", "err", err)
		return
	}
	if acquired {
		a.becomeLeader(ctx)
	}
}

// acquire 尝试抢占领导锁（SET NX EX）。
func (a *schedulerProvider) acquire(ctx context.Context) (bool, error) {
	rdb, ok := a.universalClient()
	if !ok {
		return false, errors.New("failed to assert to UniversalClient")
	}
	return rdb.SetNX(ctx, schedulerLockKey, a.token, schedulerLockTTLSeconds*time.Second).Result()
}

// renew CAS 续约。held=false 表示锁已过期或被其它实例持有（确定易主）。
func (a *schedulerProvider) renew(ctx context.Context) (held bool, err error) {
	rdb, ok := a.universalClient()
	if !ok {
		return false, errors.New("failed to assert to UniversalClient")
	}
	n, err := renewLockScript.Run(ctx, rdb, []string{schedulerLockKey}, a.token, schedulerLockTTLSeconds).Int()
	if err != nil {
		return false, err
	}
	return n == 1, nil
}

// resign 让位（幂等、进程级终态）：置 resigned 屏障阻止后续抢锁，
// 摘除任务并按 CAS 释放锁（不删他人的锁）。选主循环在下一轮检查 resigned 后退出。
func (a *schedulerProvider) resign(ctx context.Context) {
	if !a.resigned.CompareAndSwap(false, true) {
		return
	}
	if !a.IsLeader() {
		return
	}
	rdb, ok := a.universalClient()
	if !ok {
		return
	}
	a.stopJobs(ctx)
	if _, err := releaseLockScript.Run(ctx, rdb, []string{schedulerLockKey}, a.token).Result(); err == nil {
		g.Log().Warning(ctx, "调度器：成功释放领导者锁", "token", a.token)
	}
	a.leader.Store(false)
}

func (a *schedulerProvider) becomeLeader(ctx context.Context) {
	a.leader.Store(true)
	g.Log().Warning(ctx, "调度器：当选为领导者，正在注册后台任务", "token", a.token)
	for _, job := range schedulerJobs {
		// 防御性先按名移除：防止本进程此前当选周期残留同名注册导致 AddSingleton 失败。
		gcron.Remove(job.name)
		job.start(ctx)
	}
}

func (a *schedulerProvider) loseLeadership(ctx context.Context) {
	g.Log().Warning(ctx, "调度器：失去领导者身份，正在移除后台任务", "token", a.token)
	a.leader.Store(false)
	a.stopJobs(ctx)
}

func (a *schedulerProvider) stopJobs(ctx context.Context) {
	for _, job := range schedulerJobs {
		gcron.Remove(job.name)
	}
}
