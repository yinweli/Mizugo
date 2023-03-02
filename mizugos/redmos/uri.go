package redmos

import (
	"fmt"
	"strings"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/cast"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/yinweli/Mizugo/mizugos/ctxs"
)

// RedisURI redis連接字串, 選項字串是仿造mongo配置字串做出來的, 選項字串語法如下
//   - redisdb://[username:password@]host1:port1[,host2:port2,...,hostN:portN]/[?options]
//
// 語法中各組件的說明如下
//   - redisdb://
//     @必選 必須是 'redisdb://'
//   - [username:password@]
//     @必選 指定連接時使用的帳號密碼
//   - host1:port1[,host2:port2,...,hostN:portN]/
//     @必選 指定連接位址與埠號, 如果要連接到叢集的話, 就需要設定多組位址與埠號
//   - [?options]
//     @必選 設定連接選項, 選項以'&'符號分隔, 例如: name=value&name=value
//
// 以下是連接選項說明
//   - clientName
//     設置連接名稱
//   - maxRetries
//     最大重試次數, 預設值為3, -1表示關閉此功能
//   - minRetryBackoff
//     最小重試時間間隔, 預設值為8毫秒, -1表示關閉此功能
//   - maxRetryBackoff
//     最大重試時間間隔, 預設值為512毫秒, -1表示關閉此功能
//   - dialTimeout
//     連接超時時間, 預設值為5秒
//   - readTimeout
//     socket讀取超時時間, 超時會導致命令失敗, 預設值為3秒, -1表示關閉此功能
//   - writeTimeout
//     socket寫入超時時間, 超時會導致命令失敗, 預設值為3秒, -1表示關閉此功能
//   - contextTimeoutEnabled
//     控制客戶端是否遵守上下文超時和截止日期
//   - poolFIFO
//     設置連接池的類型, true表示FIFO池, false表示LIFO池
//     與LIFO相比, FIFO的開銷略高, 但它有助於更快地關閉空閒連接, 從而減少池大小
//   - poolSize
//     最大socket連接數, 預設值為CPU數量*10
//   - poolTimeout
//     連接池超時時間, 超時會導致從連接池中獲得連接失敗, 預設值為readTimeout + 1秒
//   - minIdleConns
//     連接池中最小的閒置連接數量, 請注意建立新連接是很慢的
//   - maxIdleConns
//     連接池中最大的閒置連接數量
//   - connMaxIdleTime
//     連接閒置時間, 連接如果閒置超過此時間將會被刪除, 應該小於服務器的超時時間, 預設值為5分鐘, -1表示關閉此功能
//   - connMaxLifetime
//     連接生存時間, 連接如果超過此時間將會被刪除, 預設值為不刪除
//   - maxRedirects
//     @叢集專用 當重定向時的最大重試次數, 預設值為8次
//   - readOnly
//     @叢集專用 控制是否在從屬節點上啟用只讀命令
//   - routeByLatency
//     @叢集專用 控制是否允許將只讀命令路由到最近的主節點或從節點, 它會自動啟用readOnly
//   - routeRandomly
//     @叢集專用 控制是否允許將只讀命令路由到隨機主節點或從節點, 它會自動啟用readOnly
//
// 也可以到以下網址查看選項詳細說明
//   - https://github.com/redis/go-redis/blob/master/options.go
type RedisURI string

// Connect 連接到資料庫
func (this RedisURI) Connect(ctx ctxs.Ctx) (client redis.UniversalClient, err error) {
	option, err := this.option()

	if err != nil {
		return nil, fmt.Errorf("redisURI start: %w", err)
	} // if

	client = redis.NewUniversalClient(option)

	if _, err = client.Ping(ctx.Ctx()).Result(); err != nil {
		return nil, fmt.Errorf("redisURI start: %w", err)
	} // if

	return client, nil
}

// option 取得連接選項
func (this RedisURI) option() (option *redis.UniversalOptions, err error) {
	uri := string(this)
	option = &redis.UniversalOptions{}

	// 檢驗前置詞
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
			} // switch
		} // for
	} // if

	return option, nil
}

// MongoURI mongo連接字串, 選項字串語法如下
//   - mongodb://[username:password@]host1[:port1][,host2[:port2],...,hostN[:portN]][?options]]
//
// 語法中各組件的說明如下
//   - mongodb://
//     @必選 必須是 'mongodb://'
//   - [username:password@]
//     @可選 指定連接時使用的帳號密碼
//   - host1[:port1][,host2[:port2],...,hostN[:portN]]
//     @必選 指定連接位址與埠號, 如果要連接到分片叢集的話, 就需要設定多組位址與埠號
//     如果省略埠號的話, 就會使用mongo預設的埠號27017
//   - [?options]
//     @可選 設定連接選項, 選項以'&'符號分隔, 例如: name=value&name=value
//
// 以下是連接選項說明, 選項中有關於時間值的單位都是毫秒
//   - connectTimeoutMS
//     連接超時時間, 預設值為30000
//   - timeoutMS
//     命令超時時間, 預設值為0, 0表示不會超時
//   - maxIdleTimeMS
//     連接閒置時間, 連接如果閒置超過此時間將會被刪除, 預設值為0, 0表示不刪除
//   - heartbeatFrequencyMS
//     心跳檢查時間, 預設值為10000
//   - socketTimeoutMS
//     socket操作超時時間, 超時會導致失敗, 預設值為0, 0表示不會超時
//   - serverSelectionTimeoutMS
//     尋找伺服器超時時間
//   - minPoolSize
//     最小連接池大小, 預設值為0
//   - maxPoolSize
//     最大連接池大小, 預設值為100
//   - replicaSet
//     集群的副本集名稱, 副本集中的所有節點必須具有相同的副本集名稱, 否則客戶端不會將它們視為副本集的一部分
//
// 也可以到以下網址查看選項詳細說明
//   - https://www.mongodb.com/docs/drivers/go/current/fundamentals/connection/#std-label-golang-connection-options
type MongoURI string

// Connect 連接到資料庫
func (this MongoURI) Connect(ctx ctxs.Ctx) (client *mongo.Client, err error) {
	option, err := this.option()

	if err != nil {
		return nil, fmt.Errorf("mongoURI start: %w", err)
	} // if

	client, err = mongo.Connect(ctx.Ctx(), option)

	if err != nil {
		return nil, fmt.Errorf("mongoURI start: %w", err)
	} // if

	if err = client.Ping(ctx.Ctx(), readpref.Primary()); err != nil {
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
