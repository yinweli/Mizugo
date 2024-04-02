![license](https://img.shields.io/github/license/yinweli/Mizugo)
![lint](https://github.com/yinweli/Mizugo/actions/workflows/lint.yml/badge.svg)
![test](https://github.com/yinweli/Mizugo/actions/workflows/test.yml/badge.svg)
![codecov](https://codecov.io/gh/yinweli/Mizugo/branch/main/graph/badge.svg?token=1DGCDV1S69)

# Mizugo
以[go]做成的遊戲伺服器框架, 包括TCP網路通訊, 資料庫組件等

# 分支列表
| 分支                | 說明                                                     |
|:--------------------|:---------------------------------------------------------|
| main                | 主分支                                                   |
| client-unity        | 客戶端組件分支, 提供給[unity]的[package manager]安裝用   |
| client-unity-sample | 客戶端組件範例分支                                       |
| proto-unity         | protobuf組件分支, 提供給[unity]的[package manager]安裝用 |

# 系統需求
* [go]1.19以上
* [proto]3以上

# 如何安裝
* 安裝[go]
* 安裝[protoc]
* 安裝[protoc-go], 在終端執行以下命令
  ```sh
  go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
  ```
* 安裝[mizugo], 在終端執行以下命令
  ```sh
  go get github.com/yinweli/Mizugo
  ```

# 管理器
[mizugo]實際上是多個伺服器工具的集合, 其中最重要的功能都實作為管理器, 使用者可以依自己的需求來建立並使用這些管理器  
以下表格概述了[mizugo]中各個管理器的基本資訊, 包括它們的功能、程式碼位置、軟體包與類別名稱, 以及如何創建這些管理器的方法

| 名稱           | 程式碼位置     | 軟體包與類別名稱   | 建立函式              |
|:---------------|:---------------|:-------------------|:----------------------|
| 配置管理器     | mizugo/configs | configs.Configmgr  | configs.NewConfigmgr  |
| 實體管理器     | mizugo/entitys | entitys.Entitymgr  | entitys.NewEntitymgr  |
| 標籤管理器     | mizugo/labels  | labels.Labelmgr    | labels.NewLabelmgr    |
| 日誌管理器     | mizugo/loggers | loggers.Logmgr     | loggers.NewLogmgr     |
| 度量管理器     | mizugo/metrics | metrics.Metricsmgr | metrics.NewMetricsmgr |
| 網路管理器     | mizugo/nets    | nets.Netmgr        | nets.NewNetmgr        |
| 執行緒池管理器 | mizugo/pools   | pools.Poolmgr      | pools.DefaultPool     |
| 資料庫管理器   | mizugo/redmos  | redmos.Redmomgr    | redmos.NewRedmomgr    |

需要注意的是, 執行緒池管理器與其他管理器略有不同, 雖然它確實有一個 pools.NewPoolmgr 函式, 但建議使用時採用事先建立好的 pools.DefaultPool 實例

# 伺服器範例
伺服器的範例位於 support/test-server 路徑下, 以下介紹一些關鍵的目錄和檔案及其用途  
- features: 此目錄存放各種管理器以及它們的初始化和結束處理函式
- entrys: 此目錄存放入口組件以及它們的初始化和結束處理函式. 每一個入口都會監聽一個特定的埠號, 以便當客戶端通過該埠號建立連接時, 能夠使用預設的模組來創建相應的實例
- modules: 此目錄存放負責處理訊息交換的模組, 這些模組實現了伺服器與客戶端間通訊的具體邏輯
- querys: 此目錄存放負責處理與資料庫的組件, 包括數據的新增與查詢等操作
- cmd/server.go: 伺服器啟動的入口點, 這個檔案中的 main 函式負責錯誤處理、初始化管理器、初始化入口以及保持伺服器的持續運行

# 如何使用客戶端組件
請參閱[客戶端組件說明](support/client-unity/Packages/com.fouridstudio.mizugo-client-unity/README.md)  
請參閱[proto組件說明](support/client-unity/Packages/com.fouridstudio.mizugo-proto-unity/README.md)

# 專案目錄說明
| 目錄                   | 說明                            |
|:-----------------------|:--------------------------------|
| mizugo                 | mizugo程式碼                    |
| mizugo/configs         | 配置組件                        |
| mizugo/cryptos         | 加密/解密組件                   |
| mizugo/ctxs            | context組件                     |
| mizugo/entitys         | 實體, 模組, 事件組件            |
| mizugo/helps           | 協助組件                        |
| mizugo/iaps            | 購買驗證組件                    |
| mizugo/labels          | 標籤組件                        |
| mizugo/loggers         | 日誌組件                        |
| mizugo/metrics         | 度量組件                        |
| mizugo/msgs            | 封包結構                        |
| mizugo/nets            | 網路組件                        |
| mizugo/pools           | 執行緒池組件                    |
| mizugo/procs           | 處理器組件                      |
| mizugo/redmos          | 雙層式資料庫組件(redis + mongo) |
| support                | 支援專案                        |
| support/client-unity   | unity客戶端組件                 |
| support/proto          | proto定義檔                     |
| support/proto/mizugo   | 內部proto定義檔                 |
| support/proto/test     | 測試proto定義檔                 |
| support/test-client-cs | unity測試客戶端                 |
| support/test-client-go | go測試客戶端                    |
| support/test-server    | 測試伺服器                      |
| testdata               | 測試資料與測試工具              |

# 軟體包階層
| 軟體包名稱列表                                     |
|:---------------------------------------------------|
| testdata                                           |
| ctxs, helps, msgs                                  |
| cryptos, iaps, nets, pools, procs                  |
| configs, entitys, labels, loggers, metrics, redmos |

下面的軟體包可以引用上面的軟體包  
上面的軟體包不能引用下面的軟體包  
相同階層的不能互相引用

# Task命令說明
輸入 `task 命令名稱` 來執行命令, 如果無法使用, 表示還沒有安裝[task]

| 命令名稱       | 命令說明         |
|:---------------|:-----------------|
| lint           | 進行程式碼檢查   |
| test           | 進行程式碼測試   |
| bench          | 進行效能測試     |
| proto          | 更新訊息         |
| subtree        | 更新子專案分支   |
| stop           | 停止容器         |
| db             | 啟動資料庫       |

# JetBrains licenses
[mizugo]使用了JetBrains的Goland的免費開發許可, 在此表示感謝  
<img src="https://resources.jetbrains.com/storage/products/company/brand/logos/jb_beam.png" alt="JetBrains Logo (Main) logo." style="width:200px;">
<img src="https://resources.jetbrains.com/storage/products/company/brand/logos/GoLand_icon.png" alt="GoLand logo." style="width:200px;">

[go]: https://go.dev/dl/
[mizugo]: https://github.com/yinweli/mizugo
[package manager]: https://docs.unity3d.com/Manual/Packages.html
[proto]: https://github.com/protocolbuffers/protobuf
[protoc-go]: https://github.com/protocolbuffers/protobuf-go
[protoc]: https://github.com/protocolbuffers/protobuf
[task]: https://taskfile.dev/
[unity]: https://unity.com/