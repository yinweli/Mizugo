package redmos

import (
	"context"
	"fmt"
	"strings"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/cast"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// RedisURI Redis 連線 URI, 選項語法刻意貼近 Mongo 的連線字串風格, 基本格式:
//
//	redisdb://[username:password@]host1:port1[,host2:port2,...,hostN:portN]/[?options]
//
// 組成說明:
//   - redisdb://
//     @必選 必須是 'redisdb://'
//   - [username:password@]
//     @可選 連線認證的帳號/密碼, 例如: user:pass@
//   - host1:port1[,host2:port2,...,hostN:portN]/
//     @必選 至少一組 host:port, 若為叢集或哨兵, 可填多組位址, 以逗點分隔
//   - [?options]
//     @可選 連線參數, 以 & 串接, 例如: name=value&name=value
//
// 支援的 options(對應 go-redis UniversalOptions):
//   - clientName              : 連線名稱
//   - dbid                    : 資料庫編號
//   - maxRetries              : 最大重試次數(預設 3; -1 表示停用)
//   - minRetryBackoff         : 最小重試間隔(預設 8ms; -1 表示停用)
//   - maxRetryBackoff         : 最大重試間隔(預設 512ms; -1 表示停用)
//   - dialTimeout             : 撥號逾時(預設 5s)
//   - readTimeout             : 讀取逾時(預設 3s; -1 表示停用)
//   - writeTimeout            : 寫入逾時(預設 3s; -1 表示停用)
//   - contextTimeoutEnabled   : 是否遵守 context 的逾時/截止
//   - poolFIFO                : 連線池佇列策略(true=FIFO; false=LIFO)
//   - poolSize                : 最大連線數(預設 CPU*10)
//   - poolTimeout             : 自連線池取連線的逾時(預設 readTimeout+1s)
//   - minIdleConns            : 最小閒置連線數
//   - maxIdleConns            : 最大閒置連線數
//   - connMaxIdleTime         : 連線最大閒置時間(預設 5m; -1 表示停用)
//   - connMaxLifetime         : 連線最長存活時間(預設不限制)
//   - maxRedirects            : @叢集用; 重導向時最大重試次數(預設 8)
//   - readOnly                : @叢集用; 是否允許在從節點執行唯讀命令
//   - routeByLatency          : @叢集用; 唯讀命令依延遲路由(會自動啟用 readOnly)
//   - routeRandomly           : @叢集用; 唯讀命令隨機路由(會自動啟用 readOnly)
//   - masterName              : @哨兵用; 主節點名稱(指定後採哨兵模式連線)
//
// 參考:  https://github.com/redis/go-redis/blob/master/options.go
type RedisURI string

// Connect 連接到資料庫
func (this RedisURI) Connect(ctx context.Context) (client redis.UniversalClient, err error) {
	option, err := this.option()

	if err != nil {
		return nil, fmt.Errorf("redisURI start: %w", err)
	} // if

	client = redis.NewUniversalClient(option)

	if _, err = client.Ping(ctx).Result(); err != nil {
		return nil, fmt.Errorf("redisURI start: %w", err)
	} // if

	return client, nil
}

// option 取得連接選項
func (this RedisURI) option() (option *redis.UniversalOptions, err error) {
	uri := string(this)
	option = &redis.UniversalOptions{}

	// 檢查前綴
	if prefix := "redisdb://"; strings.HasPrefix(uri, prefix) {
		uri = strings.TrimPrefix(uri, prefix)
	} else {
		return nil, fmt.Errorf("redisURI option: prefix must be \"%v\"", prefix)
	} // if

	// 取得帳號密碼
	if index := strings.Index(uri, "@"); index != -1 {
		block := uri[:index]
		uri = uri[index+1:]

		if block != "" {
			if index = strings.Index(block, ":"); index != -1 {
				option.Username = block[:index]
				option.Password = block[index+1:]
			} else {
				return nil, fmt.Errorf("redisURI option: invalid username/password")
			} // if
		} // if
	} // if

	// 取得位址埠號
	if index := strings.Index(uri, "/"); index != -1 {
		block := uri[:index]
		uri = uri[index+1:]

		if block != "" {
			for _, itor := range strings.Split(block, ",") {
				if index = strings.Index(itor, ":"); index != -1 {
					option.Addrs = append(option.Addrs, itor)
				} // if
			} // for
		} // if
	} // if

	if len(option.Addrs) == 0 {
		return nil, fmt.Errorf("redisURI option: host not found")
	} // if

	// 取得選項
	if uri != "" {
		if prefix := "?"; strings.HasPrefix(uri, prefix) {
			uri = strings.TrimPrefix(uri, prefix)
		} else {
			return nil, fmt.Errorf("redisURI option: option must start by \"?\"")
		} // if

		for _, itor := range strings.Split(uri, "&") {
			name, value, found := strings.Cut(itor, "=")

			if found == false {
				return nil, fmt.Errorf("redisURI option: invalid option format: %v", itor)
			} // if

			switch name {
			case "clientName":
				option.ClientName = cast.ToString(value)

			case "dbid":
				option.DB = cast.ToInt(value)

			case "maxRetries":
				option.MaxRetries = cast.ToInt(value)

			case "minRetryBackoff":
				option.MinRetryBackoff = cast.ToDuration(value)

			case "maxRetryBackoff":
				option.MaxRetryBackoff = cast.ToDuration(value)

			case "dialTimeout":
				option.DialTimeout = cast.ToDuration(value)

			case "readTimeout":
				option.ReadTimeout = cast.ToDuration(value)

			case "writeTimeout":
				option.WriteTimeout = cast.ToDuration(value)

			case "contextTimeoutEnabled":
				option.ContextTimeoutEnabled = cast.ToBool(value)

			case "poolFIFO":
				option.PoolFIFO = cast.ToBool(value)

			case "poolSize":
				option.PoolSize = cast.ToInt(value)

			case "poolTimeout":
				option.PoolTimeout = cast.ToDuration(value)

			case "minIdleConns":
				option.MinIdleConns = cast.ToInt(value)

			case "maxIdleConns":
				option.MaxIdleConns = cast.ToInt(value)

			case "connMaxIdleTime":
				option.ConnMaxIdleTime = cast.ToDuration(value)

			case "connMaxLifetime":
				option.ConnMaxLifetime = cast.ToDuration(value)

			case "maxRedirects":
				option.MaxRedirects = cast.ToInt(value)

			case "readOnly":
				option.ReadOnly = cast.ToBool(value)

			case "routeByLatency":
				option.RouteByLatency = cast.ToBool(value)

			case "routeRandomly":
				option.RouteRandomly = cast.ToBool(value)

			case "masterName":
				option.MasterName = value
			} // switch
		} // for
	} // if

	return option, nil
}

// add 新增連接選項
func (this RedisURI) add(option string) RedisURI {
	uri := string(this)

	if strings.Contains(uri, "?") {
		uri += "&"
	} else {
		uri += "?"
	} // if

	return RedisURI(uri + option)
}

// MongoURI Mongo 連線 URI, 基本格式:
//
//	mongodb://[username:password@]host1[:port1][,host2[:port2],...,hostN[:portN]]/[?options]]
//
// 組成說明:
//   - mongodb://
//     @必選 必須是 'mongodb://'
//   - [username:password@]
//     @可選 連線認證的帳號/密碼, 例如: user:pass@
//   - host1[:port1][,host2[:port2],...,hostN[:portN]]
//     @必選 至少一組 host:port, 可填多組以支援分片/副本集; 若省略 port, 預設為 27017
//   - [?options]
//     @可選 連線參數, 以 & 串接, 例如: name=value&name=value
//
// 支援的 options (時間值單位為毫秒):
//   - connectTimeoutMS           : 連線逾時(預設 30000)
//   - timeoutMS                  : 命令逾時(預設 0; 0 表示不逾時)
//   - maxIdleTimeMS              : 連線最大閒置時間(預設 0; 0 表示不刪除)
//   - heartbeatFrequencyMS       : 心跳頻率(預設 10000)
//   - socketTimeoutMS            : Socket 操作逾時(預設 0; 0 表示不逾時)
//   - serverSelectionTimeoutMS   : 伺服器選擇逾時
//   - minPoolSize                : 最小連線池大小(預設 0)
//   - maxPoolSize                : 最大連線池大小(預設 100)
//   - replicaSet                 : 副本集名稱(叢集節點需一致)
//
// 參考:  https://www.mongodb.com/docs/drivers/go/current/fundamentals/connection/#std-label-golang-connection-options
type MongoURI string

// Connect 連接到資料庫
func (this MongoURI) Connect(ctx context.Context) (client *mongo.Client, err error) {
	option, err := this.option()

	if err != nil {
		return nil, fmt.Errorf("mongoURI start: %w", err)
	} // if

	client, err = mongo.Connect(ctx, option)

	if err != nil {
		return nil, fmt.Errorf("mongoURI start: %w", err)
	} // if

	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, fmt.Errorf("mongoURI start: %w", err)
	} // if

	return client, nil
}

// option 取得連接選項
func (this MongoURI) option() (option *options.ClientOptions, err error) {
	uri := string(this)
	option = options.Client().ApplyURI(uri)

	if err = option.Validate(); err != nil {
		return nil, fmt.Errorf("mongoURI option: %w", err)
	} // if

	return option, nil
}
