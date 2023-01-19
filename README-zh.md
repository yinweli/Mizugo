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

# 如何使用

# 專案目錄說明

| 目錄                     | 說明            |
|:-------------------------|:----------------|
| mizugos                  | 核心程式碼      |
| mizugos/configs          | 配置組件        |
| mizugos/contexts         | context組件     |
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
| support/clientcs         | cs客戶端組件    |
| support/example_clientcs | cs測試客戶端    |
| support/example_clientgo | go測試客戶端    |
| support/example_server   | 測試伺服器      |
| support/proto            | proto定義檔     |
| support/proto/example    | 測試proto定義檔 |
| support/proto/mizugo     | 內部proto定義檔 |

[go]: https://go.dev/dl/
[proto]: https://github.com/protocolbuffers/protobuf
[protoc-go]: https://github.com/protocolbuffers/protobuf-go
[protoc]: https://github.com/protocolbuffers/protobuf