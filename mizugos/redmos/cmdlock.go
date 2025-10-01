package redmos

import (
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Lock 鎖定行為
//
// 以索引鍵(Key)在主要資料庫執行分布式鎖定
//
// 事前準備:
//   - 設定 Key: 不可為空字串
//   - 設定 Token: 不可為空字串
//
// 注意:
//   - 本行為僅使用主要資料庫, 次要資料庫不參與
//   - 鎖定後需搭配 Unlock 行為解鎖
//   - Token 必須為唯一識別字串; Unlock 會核對 Token, 避免誤解他人持有的鎖
type Lock struct {
	Behave                  // 行為物件
	Key    string           // 索引值
	Token  string           // 識別字串
	ttl    time.Duration    // 逾時時間(供測試調整)
	cmd    *redis.StatusCmd // 命令結果
}

// Prepare 前置處理
func (this *Lock) Prepare() error {
	if this.Key == "" {
		return fmt.Errorf("lock prepare: key empty")
	} // if

	if this.Token == "" {
		return fmt.Errorf("lock prepare: token empty")
	} // if

	key := fmt.Sprintf(lockFormat, this.Key)
	this.cmd = this.Major().SetArgs(this.Ctx(), key, this.Token, redis.SetArgs{Mode: "NX", TTL: this.ttl})
	return nil
}

// Complete 完成處理
func (this *Lock) Complete() error {
	ok, err := this.cmd.Result()

	if err != nil {
		return fmt.Errorf("lock complete: %w: %v", err, this.Key)
	} // if

	if ok != RedisOk {
		return fmt.Errorf("lock complete: already lock: %v", this.Key)
	} // if

	return nil
}

// Unlock 解鎖行為
//
// 以索引鍵(Key)在主要資料庫執行解鎖
//
// 事前準備:
//   - 設定 Key: 不可為空字串
//   - 設定 Token: 不可為空字串; 須與 Lock 時所用 Token 一致
//
// 注意:
//   - 本行為僅使用主要資料庫, 次要資料庫不參與
//   - 僅在 Token 相符時才會刪除鎖鍵; Token 不符或鎖已逾期/不存在時, 不會刪除
//   - 若鎖已因 ttl 逾期而自動釋放, 此操作將不會刪除任何鍵
type Unlock struct {
	Behave            // 行為物件
	Key    string     // 索引值
	Token  string     // 識別字串
	cmd    *redis.Cmd // 命令結果
}

// Prepare 前置處理
func (this *Unlock) Prepare() error {
	if this.Key == "" {
		return fmt.Errorf("unlock prepare: key empty")
	} // if

	if this.Token == "" {
		return fmt.Errorf("unlock prepare: token empty")
	} // if

	key := fmt.Sprintf(lockFormat, this.Key)
	this.cmd = this.Major().Eval(this.Ctx(), lockLUA, []string{key}, this.Token)
	return nil
}

// Complete 完成處理
func (this *Unlock) Complete() error {
	if _, err := this.cmd.Result(); err != nil {
		return fmt.Errorf("unlock complete: %w: %v", err, this.Key)
	} // if

	return nil
}

const lockFormat = "lock:%v"                                                                                        // 鎖定/解鎖索引格式
const lockLUA = `if redis.call("GET", KEYS[1]) == ARGV[1] then return redis.call("DEL", KEYS[1]) else return 0 end` // 解鎖 LUA 腳本
