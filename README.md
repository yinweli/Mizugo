[![license](https://img.shields.io/github/license/yinweli/Mizugo)](https://github.com/yinweli/Mizugo/blob/main/LICENSE)
[![lint](https://github.com/yinweli/Mizugo/actions/workflows/lint.yml/badge.svg)](https://github.com/yinweli/Mizugo/actions/workflows/lint.yml)
[![test](https://github.com/yinweli/Mizugo/actions/workflows/test.yml/badge.svg)](https://github.com/yinweli/Mizugo/actions/workflows/test.yml)
[![codecov](https://codecov.io/gh/yinweli/Mizugo/branch/main/graph/badge.svg?token=LK5HL58LSN)](https://codecov.io/gh/yinweli/Mizugo)

# Mizugo
[Mizugo] 是一款使用 [Go] 語言開發的遊戲伺服器框架, 包括 TCP 網路, 通訊協議, 資料庫組件等

# 分支列表
伺服器專案中的主要 Git 分支說明如下

| 分支                | 說明                                                         |
|:--------------------|:-------------------------------------------------------------|
| main                | 主分支                                                       |
| client-unity        | 客戶端組件分支, 提供給 [Unity] 的 [Package Manager] 安裝用   |
| client-unity-sample | 客戶端組件範例分支                                           |
| proto-unity         | protobuf組件分支, 提供給 [Unity] 的 [Package Manager] 安裝用 |

# 系統需求
- [Go] 1.24+
- [Proto] v3

# 如何安裝
- 確保在您的系統上安裝了 [Go] 語言
- 確保在您的系統上安裝了 [Protoc]
- 安裝 [Protoc-Go], 請在終端機執行以下命令
  ```sh
  go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
  ```
- 在終端機中執行以下命令以安裝最新版本
  ```sh
  go get github.com/yinweli/Mizugo/v2
  ```
- 如果需要安裝 v1 的 [Mizugo], 可以使用以下命令安裝
  ```sh
  go get github.com/yinweli/Mizugo
  ```

# 伺服器組件
下面是一個簡化的伺服器程式範例, 如果需要更詳細的範例, 請造訪 `support/test-server/` 目錄
```go
func main() {
    defer func() {
        if cause := recover(); cause != nil {
            // 處理崩潰錯誤
        } // if
    }()

    mizugos.Start() // 啟動伺服器
    ctx := ctxs.Get().WithCancel()

    // 使用者自訂的初始化程序
    // 如果有任何失敗, 執行 mizugos.Stop() 後退出

    for range ctx.Done() { // 進入無限迴圈直到執行 ctx.Cancel()
    } // for

    // 使用者自訂的結束程序
    // 如果有任何失敗, 執行 mizugos.Stop() 後退出

    mizugos.Stop() // 關閉伺服器
}
```

# 客戶端組件
請參閱 [客戶端組件說明](support/client-unity/Packages/com.fouridstudio.mizugo-client-unity/README.md)

請參閱 [Proto組件說明](support/client-unity/Packages/com.fouridstudio.mizugo-proto-unity/README.md)

# 管理器
[Mizugo] 擁有多種功能, 這些功能透過不同的管理器提供, 使用者可以根據自己的需求選擇使用這些管理器

以下提供了各管理器的基本資訊, 包括它們的軟體包與類別名稱, 以及所在的位置

| 名稱           | 軟體包與類別名稱    | 管理器位置     |
|:---------------|:--------------------|:---------------|
| 配置管理器     | configs.Configmgr   | mizugos.Config |
| 日誌管理器     | loggers.Logmgr      | mizugos.Log    |
| 網路管理器     | nets.Netmgr         | mizugos.Net    |
| 資料庫管理器   | redmos.Redmomgr     | mizugos.Redmo  |
| 實體管理器     | entitys.Entitymgr   | mizugos.Entity |
| 執行緒池管理器 | pools.Poolmgr       | mizugos.Pool   |

# 訊息處理器
[Mizugo] 提供多種訊息處理器, 使用者可以根據自己的需求選擇使用這些訊息處理器

以下提供了各訊息處理器的基本資訊, 包括它們的軟體包與類別名稱

| 名稱             | 軟體包與類別名稱 |
|:-----------------|:-----------------|
| JSON 訊息處理器  | procs.Json       |
| Proto 訊息處理器 | procs.Proto      |
| Raven 訊息處理器 | procs.Raven      |

# 專案目錄說明
| 目錄                   | 說明                 |
|:-----------------------|:---------------------|
| mizugos                | Mizugo 程式碼        |
| mizugos/configs        | 配置組件             |
| mizugos/cryptos        | 加密/解密組件        |
| mizugos/ctxs           | context 組件         |
| mizugos/entitys        | 實體, 模組, 事件組件 |
| mizugos/helps          | 協助組件             |
| mizugos/iaps           | 購買驗證組件         |
| mizugos/loggers        | 日誌組件             |
| mizugos/msgs           | 封包結構             |
| mizugos/nets           | 網路組件             |
| mizugos/pools          | 執行緒池組件         |
| mizugos/procs          | 訊息處理器組件       |
| mizugos/redmos         | 雙層式資料庫組件     |
| mizugos/trials         | 單元測試組件         |
| support                | 支援專案             |
| support/client-unity   | Unity 客戶端組件     |
| support/proto-mizugo   | 內部 Proto 定義      |
| support/proto-test     | 測試 Proto 定義      |
| support/test-client-cs | Unity 測試客戶端     |
| support/test-client-go | Go 測試客戶端        |
| support/test-server    | 測試伺服器           |
| testdata               | 測試資料             |

# 軟體包階層
| 階層名稱 | 軟體包名稱      |
|:---------|:----------------|
| 測試層   | testdata        |
| 測試層   | mizugos/trials  |
| 工具層   | mizugos/ctxs    |
| 工具層   | mizugos/helps   |
| 工具層   | mizugos/msgs    |
| 通用層   | mizugos/cryptos |
| 通用層   | mizugos/iaps    |
| 通用層   | mizugos/nets    |
| 通用層   | mizugos/pools   |
| 通用層   | mizugos/procs   |
| 組件層   | mizugos/configs |
| 組件層   | mizugos/entitys |
| 組件層   | mizugos/loggers |
| 組件層   | mizugos/redmos  |

下面階層的軟體包可以引用上面階層的軟體包

上面階層的軟體包不能引用下面階層的軟體包

相同階層的不能互相引用

# Task命令說明
輸入 `task 命令名稱` 來執行命令, 如果無法使用, 表示還沒有安裝 [Task] (安裝方式請見其官方文件)

| 命令名稱       | 命令說明         |
|:---------------|:-----------------|
| lint           | 進行程式碼檢查   |
| test           | 進行程式碼測試   |
| bench          | 進行效能測試     |
| proto          | 更新訊息         |
| subtree        | 更新子專案分支   |
| stop           | 停止容器         |
| db             | 啟動資料庫       |
| install        | 安裝工具         |

# JetBrains 感謝與授權標示
[Mizugo] 使用了 JetBrains 的 Goland 的免費開發許可, 在此表示感謝

<img src="https://resources.jetbrains.com/storage/products/company/brand/logos/jb_beam.png" alt="JetBrains Logo (Main) logo." style="width:200px;">
<img src="https://resources.jetbrains.com/storage/products/company/brand/logos/GoLand_icon.png" alt="GoLand logo." style="width:200px;">

[Go]: https://go.dev/dl/
[Mizugo]: https://github.com/yinweli/mizugo
[Package Manager]: https://docs.unity3d.com/Manual/Packages.html
[Proto]: https://github.com/protocolbuffers/protobuf
[Protoc-Go]: https://github.com/protocolbuffers/protobuf-go
[Protoc]: https://github.com/protocolbuffers/protobuf
[Task]: https://taskfile.dev/
[Unity]: https://unity.com/