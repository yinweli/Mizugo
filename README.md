![license](https://img.shields.io/github/license/yinweli/Mizugo)
![lint](https://github.com/yinweli/Mizugo/actions/workflows/lint.yml/badge.svg)
![test](https://github.com/yinweli/Mizugo/actions/workflows/test.yml/badge.svg)
![codecov](https://codecov.io/gh/yinweli/Mizugo/branch/main/graph/badge.svg?token=1DGCDV1S69)

# Mizugo
以[go]做成的遊戲伺服器框架, 包括TCP網路通訊, 資料庫組件等  

# 系統需求
* [go]1.18以上
* [proto]3以上

# 安裝說明
* 安裝[go]
* 安裝[protoc]
* 安裝[protoc-go], 在終端執行以下命令
  ```sh
  go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
  ```

# 如何使用伺服器組件

# 如何使用客戶端組件
請參閱[客戶端組件說明][client-unity]

# 專案目錄說明

| 目錄                     | 說明            |
|:-------------------------|:----------------|
| mizugos                  | 核心程式碼      |
| mizugos/configs          | 配置組件        |
| mizugos/contexts         | context組件     |
| mizugos/cryptos          | 加密/解密組件   |
| mizugos/entitys          | 實體與模組組件  |
| mizugos/errors           | 錯誤組件        |
| mizugos/events           | 事件組件        |
| mizugos/labels           | 標籤組件        |
| mizugos/logs             | 日誌組件        |
| mizugos/metrics          | 度量組件        |
| mizugos/msgs             | 封包結構        |
| mizugos/nets             | 網路組件        |
| mizugos/pools            | 執行緒池組件    |
| mizugos/procs            | 處理器組件      |
| mizugos/utils            | 協助組件        |
| support                  | 支援專案        |
| support/client-unity     | unity客戶端組件 |
| support/test-client-cs   | unity測試客戶端 |
| support/test-client-go   | go測試客戶端    |
| support/test-server      | 測試伺服器      |
| support/proto            | proto定義檔     |
| support/proto/mizugo     | 內部proto定義檔 |
| support/proto/test       | 測試proto定義檔 |

# 專案分支說明

| 分支                | 說明                                                   |
|:--------------------|:-------------------------------------------------------|
| main                | 主分支                                                 |
| client-unity        | 客戶端組件分支, 提供給[unity]的package manager安裝用   |
| client-unity-sample | 客戶端組件範例分支                                     |
| proto-unity         | protobuf組件分支, 提供給[unity]的package manager安裝用 |

# Taskfile命令速查
* lint: 進行程式碼檢查
* test: 進行程式碼測試
* bench: 進行效能測試
* proto: 產生proto message
* clientunity: 更新客戶端組件分支

[go]: https://go.dev/dl/
[proto]: https://github.com/protocolbuffers/protobuf
[protoc-go]: https://github.com/protocolbuffers/protobuf-go
[protoc]: https://github.com/protocolbuffers/protobuf
[unity]: https://unity.com/

[client-unity]: support/client-unity/Packages/com.fouridstudio.mizugo-client-unity/README.md