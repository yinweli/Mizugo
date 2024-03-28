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

# 如何使用伺服器組件
[mizugo]實際上是多個伺服器工具的集合, 最簡單啟動mizugo伺服器的程式碼範例如下  
```go
func main() {
    defer func() {
        if cause := recover(); cause != nil {
            // 處理崩潰錯誤
        } // if
    }()
    
    ctx := ctxs.Get().WithCancel()
    name := "伺服器名稱"
    mizugos.Start() // 啟動伺服器
    
    // 使用者自訂的初始化程序
    // 如果有任何失敗, 執行 mizugos.Stop() 後退出
    
    fmt.Printf("%v start\n", name)
    
    for range ctx.Done() { // 進入無限迴圈直到執行 ctx.Cancel()
    } // for
    
    // 使用者自訂的結束程序
    // 如果有任何失敗, 執行 mizugos.Stop() 後退出
    
    mizugos.Stop() // 關閉伺服器
    fmt.Printf("%v shutdown\n", name)
}
```
在 mizugos/mizugo.go 中包含了啟動/關閉伺服器的函式, 以及事先建立好的各管理器以及其取得函式  

| 項目                 | 說明           |
|:---------------------|:---------------|
| mizugos.Configmgr()  | 配置管理器     |
| mizugos.Metricsmgr() | 度量管理器     |
| mizugos.Logmgr()     | 日誌管理器     |
| mizugos.Netmgr()     | 網路管理器     |
| mizugos.Redmomgr()   | 資料庫管理器   |
| mizugos.Entitymgr()  | 實體管理器     |
| mizugos.Labelmgr()   | 標籤管理器     |
| mizugos.Poolmgr()    | 執行緒池管理器 |

這些管理器在啟動伺服器之後可用, 若在啟動伺服器之前(或是關閉伺服器之後)使用會碰到panic  
伺服器程式碼的範例可以到 support/test-server 查看  

# 如何使用客戶端組件
請參閱[客戶端組件說明](support/client-unity/Packages/com.fouridstudio.mizugo-client-unity/README.md)  
請參閱[proto組件說明](support/client-unity/Packages/com.fouridstudio.mizugo-proto-unity/README.md)

# 專案目錄說明
| 目錄                   | 說明                            |
|:-----------------------|:--------------------------------|
| mizugos                | mizugo程式碼                    |
| mizugos/configs        | 配置組件                        |
| mizugos/cryptos        | 加密/解密組件                   |
| mizugos/ctxs           | context組件                     |
| mizugos/entitys        | 實體, 模組, 事件組件            |
| mizugos/helps          | 協助組件                        |
| mizugos/iaps           | 購買驗證組件                    |
| mizugos/labels         | 標籤組件                        |
| mizugos/loggers        | 日誌組件                        |
| mizugos/metrics        | 度量組件                        |
| mizugos/msgs           | 封包結構                        |
| mizugos/nets           | 網路組件                        |
| mizugos/pools          | 執行緒池組件                    |
| mizugos/procs          | 處理器組件                      |
| mizugos/redmos         | 雙層式資料庫組件(redis + mongo) |
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