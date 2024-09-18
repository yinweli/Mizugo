# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Planning]

## [Unrelease]

## [1.1.35] - 2024-09-18
### Added
- 新增redmos.GetEx與redmos.SetEx命令
### Changed
- 移除redmos.Get與redmos.Set中的Expire Time

## [1.1.34] - 2024-09-06
### Added
- 在Redis Get的時候加上Expire Time，並且每次Get都會重置

## [1.1.33] - 2024-09-03
### Changed
- 把EquateApproxProtoTimestamp搬移到helps軟體包中

## [1.1.32] - 2024-09-03
### Added
- 新增timestamppb.Timestamp的比對選項

## [1.1.31] - 2024-09-03
### Changed
- proto比對改用go-cmp, 並添加比對選項列表

## [1.1.30] - 2024-08-29
### Changed
- 更新lint配置
- 更新redmos說明

## [1.1.29] - 2024-08-22
### Fixed
- 接收訊息時加上長度上限檢查

## [1.1.28] - 2024-08-20
### Added
- 新增主要資料庫搜尋行為
- 新增次要資料庫搜尋行為

## [1.1.27] - 2024-08-14
### Fixed
- 修正錯誤輸出

## [1.1.26] - 2024-07-19
### Fixed
- 修正標頭類型錯誤

## [1.1.25] - 2024-07-19
### Changed
- 把封包標頭長度與封包大小上限改為可由外部控制

## [1.1.24] - 2024-07-18
### Changed
- 封包標頭長度改為4byte
- 封包大小上限改為1MB

## [1.1.23] - 2024-07-16
### Deleted
- 移除ctxs組件

## [1.1.22] - 2024-07-12
### Added
- 混合資料庫新增取得主要資料庫與次要資料庫

## [1.1.21] - 2024-07-12
### Changed
- 更新IAP組件, 讓外部有更多的重試與逾時參數可用

## [1.1.20] - 2024-07-10
### Added
- 新增ProtoEqual, MapToArray函式

## [1.1.19] - 2024-06-25
### Added
- 新增FlagszAND, FlagszOR, FlagszXOR函式

## [1.1.18] - 2024-06-20
### Changed
- 更新Fit組件說明
- StdColor新增失敗旗標機制
- StdColor新增Out, Err函式

## [1.1.17] - 2024-06-06
### Fixed
- 修正OneStore的IAP驗證錯誤

## [1.1.16] - 2024-06-05
### Fixed
- 修正OneStore的IAP驗證錯誤

## [1.1.15] - 2024-06-03
### Added
- 新增OneStore的IAP組件

## [1.1.14] - 2024-05-28
### Changed
- 處理介面新增取得訊息處理功能

## [1.1.13] - 2024-05-27
### Changed
- 變更取得回應列表的方式

## [1.1.12] - 2024-05-27
### Added
- 新增取得回應列表的方式

## [1.1.11] - 2024-05-23
### Added
- 新增CPU剖析工具

## [1.1.10] - 2024-05-22
### Changed
- 讓Raven的測試系列函式會印出錯誤內容

## [1.1.9] - 2024-05-13
### Added
- 新增ReflectFieldValue函式

## [1.1.8] - 2024-05-13
### Added
- 新增WriteFile函式

### Changed
- 將WriteTemplate改為TemplateBuild, 並且只負責將新增部分添加到字串尾端

## [1.1.7] - 2024-05-13
### Added
- 新增WriteTemplate函式

## [1.1.6] - 2024-05-08
### Changed
- 日誌的Caller函式新增可選的紀錄簡單函式名稱功能

## [1.1.5] - 2024-05-06
### Added
- 新增Before, Beforef, Beforefx函式
- 新增After, Afterf, Afterfx函式
- 新增Betweenf, Betweenfx函式
- 新增Overlapff, Overlapffx函式

## [1.1.4] - 2024-05-03
### Changed
- 保留Date函式的年/月/日參數, 只有時/分/秒/毫秒可以被省略

## [1.1.3] - 2024-05-03
### Fixed
- 修正每月到期計算錯誤

## [1.1.2] - 2024-05-03
### Added
- 新增JsonString函式
- 新增ToProtoAny函式
- 新增Timef函式

### Removed
- 移除cast工具
- 移除Beforef函式
- 移除Afterf函式
- 移除Betweenf函式
- 移除SameDay函式

### Changed
- ProtoAny函式改名為FromProtoAny
- ProtoJson函式改名為ProtoString
- Date函式變得更容易使用

## [1.1.1] - 2024-04-30
### Added
- 新增函式Beforef, Afterf

## [1.1.0] - 2024-04-26
### Added
- 新增單元測試工具集trial
- 新增信號調度組件trigger
- 新增訊息處理組件procs/raven
- 新增產生錯誤工具

### Changed
- 更新網路組件
- 更新實體組件

## [1.0.0] - 2024-04-15