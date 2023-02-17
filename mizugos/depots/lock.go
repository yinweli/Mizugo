package depots

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Lock 鎖定行為, 以索引字串到redis中執行分布式鎖定, 主要用於避免同時執行客戶端動作,
// 當操作完成後, 執行 Unlock 行為來解除鎖定
type Lock struct {
	Key    string         // 鎖定索引字串
	time   time.Duration  // 鎖定超時時間, redis的超時時間不能低於1秒
	incr   *redis.IntCmd  // 遞增命令結果
	expire *redis.BoolCmd // 超時命令結果
}

// Prepare 前置處理
func (this *Lock) Prepare(ctx context.Context, majorRunner MajorRunner, _ MinorRunner) error {
	this.incr = majorRunner.Incr(ctx, this.Key)
	this.expire = majorRunner.Expire(ctx, this.Key, this.time)
	return nil
}

// Result 結果處理
func (this *Lock) Result() error {
	value, err := this.incr.Result()

	if err != nil {
		return fmt.Errorf("lock result: %w", err)
	} // if

	if value != 1 {
		return fmt.Errorf("lock result: already lock")
	} // if

	if _, err = this.expire.Result(); err != nil {
		return fmt.Errorf("lock result: %w", err)
	} // if

	return nil
}

// Unlock 解鎖行為, 解除被 Lock 行為鎖定的索引
type Unlock struct {
	Key string        // 鎖定索引字串
	del *redis.IntCmd // 刪除命令結果
}

// Prepare 前置處理
func (this *Unlock) Prepare(ctx context.Context, majorRunner MajorRunner, _ MinorRunner) error {
	this.del = majorRunner.Del(ctx, this.Key)
	return nil
}

// Result 結果處理
func (this *Unlock) Result() error {
	if _, err := this.del.Result(); err != nil {
		return fmt.Errorf("unlock result: %w", err)
	} // if

	return nil
}
