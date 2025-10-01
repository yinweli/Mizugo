# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.4.5] - 2025-10-01
### Added
- 新增proto比對時timestamppb.Timestamp的比對選項

## [0.4.4] - 2024-08-28
### Fixed
- 傳輸流內部的錯誤都算是斷線

## [0.4.3] - 2024-07-19
### Fixed
- 修正標頭類型錯誤

## [0.4.2] - 2024-07-19
### Changed
- 把封包標頭長度與封包大小上限改為可由外部控制

## [0.4.1] - 2024-07-18
### Changed
- 封包標頭長度改為4byte
- 封包大小上限改為1MB

## [0.4.0] - 2024-04-26
### Added
- 新增訊息處理組件Raven

## [0.3.1] - 2024-04-22
### Changed
- JsonProc更名為ProcJson
- ProtoProc更名為ProcProto

## [0.3.0] - 2024-04-19
### Changed
- 更新說明文件
- 改變使用方式, 加密/解密與base64不再與JsonProc/ProtoProc組件綁定
- TCPClient新增SetEvent, SetProc, SetCodec函式, 增加編碼/解碼的使用靈活度

## [0.2.4] - 2023-12-06
### Fixed
- 修正錯誤時只需要通知錯誤事件即可

## [0.2.3] - 2023-12-05
### Removed
- 移除Host與Port函式

## [0.2.2] - 2023-08-25
### Fixed
- 修正斷線時不再拋出異常, 只會剩下斷線事件

## [0.2.1] - 2023-05-26
### Fixed
- 修正收到被切分的封包造成收包錯誤

## [0.2.0] - 2023-03-08
### Changed
- Json處理增加base64與des-cbc功能及開關
- Proto處理增加base64與des-cbc功能及開關
### Removed
- 移除PList處理器

## [0.1.5] - 2023-02-09
### Changed
- 變更客戶端組件介面(IClient)
- 更新說明文件
### Fixed
- 修正連線逾時檢測

## [0.1.4] - 2023-02-08
### Changed
- 變更接收與傳送程序, 改用NetworkStream的原始功能
### Added
- 新增RecvHeaderException異常
- 新增RecvPacketException異常
### Removed
- 移除ReceiveException異常

## [0.1.3] - 2023-02-06
### Changed
- 更新說明文件

## [0.1.2] - 2023-02-06
### Changed
- 更新說明文件

## [0.1.1] - 2023-02-06
### Changed
- 更新說明文件
- 單元測試移除不需要的訊息
- 單元測試新增循環測試

## [0.1.0] - 2023-02-04
### Added
- 新增客戶端網路組件
- 新增Json訊息處理器
- 新增Proto訊息處理器
- 新增PList訊息處理器
- 新增完整範例專案