package redmos

import (
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const lockFormat = "lock:%v" // 鎖定/解鎖索引格式

// Lock 鎖定行為, 以索引字串到redis中執行分布式鎖定, 避免同時執行客戶端動作, 使用上有以下幾點須注意
//   - 執行前設定好 Key 並且不能為空字串
//   - 鎖定完成後, 需要執行 Unlock 行為來解除鎖定
//   - 鎖定後會在 Timeout 之後自動解鎖, 避免死鎖
type Lock struct {
	Behave                // 行為物件
	Key    string         // 索引值
	time   time.Duration  // 超時時間, redis的超時時間不能低於1秒; 這個欄位是為了讓單元測試可以設定較短的超時時間
	incr   *redis.IntCmd  // 遞增命令結果
	expire *redis.BoolCmd // 超時命令結果
}

// Prepare 前置處理
func (this *Lock) Prepare() error {
	if this.Key == "" {
		return fmt.Errorf("lock prepare: key empty")
	} // if

	key := fmt.Sprintf(lockFormat, this.Key)
	this.incr = this.Major().Incr(this.Ctx(), key)
	this.expire = this.Major().Expire(this.Ctx(), key, this.time)
	return nil
}

// Complete 完成處理
func (this *Lock) Complete() error {
	data, err := this.incr.Result()

	if err != nil {
		return fmt.Errorf("lock complete: %w: %v", err, this.Key)
	} // if

	if data != 1 {
		return fmt.Errorf("lock complete: already lock: %v", this.Key)
	} // if

	if _, err = this.expire.Result(); err != nil {
		return fmt.Errorf("lock complete: %w: %v", err, this.Key)
	} // if

	return nil
}

// Unlock 解鎖行為, 解除被 Lock 行為鎖定的索引, 使用上有以下幾點須注意
//   - 執行前設定好 Key 並且不能為空字串
type Unlock struct {
	Behave               // 行為物件
	Key    string        // 索引值
	del    *redis.IntCmd // 刪除命令結果
}

// Prepare 前置處理
func (this *Unlock) Prepare() error {
	if this.Key == "" {
		return fmt.Errorf("unlock prepare: key empty")
	} // if

	key := fmt.Sprintf(lockFormat, this.Key)
	this.del = this.Major().Del(this.Ctx(), key)
	return nil
}

// Complete 完成處理
func (this *Unlock) Complete() error {
	if _, err := this.del.Result(); err != nil {
		return fmt.Errorf("unlock complete: %w: %v", err, this.Key)
	} // if

	return nil
}
