# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [2.0.6] - 2025-11-20
### Added
- StringDuration: 新增支援解析包含時間單位的字串, 轉換為 time.Duration
### Changed
- StrPercentage → StringPercentage: 統一命名規則, 回傳格式不變

## [2.0.5] - 2025-10-31
### Added
- WaitFor: 新增可等待條件函式達成或超時的功能, 支援設定逾時時間

## [2.0.4] - 2025-10-28
### Added
- LessBase58 / LessBase80: 提供字串數值排序用比較器, 可直接供 sort.Slice 或 slices.SortFunc 使用
- RankBaseN(model string): 建立進制字元表對應表
- LessBaseN: 泛用進制字串比較函式, 支援忽略前導零, 非法字元回退字典序
### Changed
- ToBaseN 由 func ToBaseN(model string, input uint64) → func ToBaseN(input uint64, model string)
- FromBaseN 由 func FromBaseN(model, input string) → func FromBaseN(input, model string, rank *[256]int)
- ToBase58 / FromBase58 / ToBase80 / FromBase80 改為使用排序表
- 新增溢位檢查, 防止 uint64 轉換時超界
- 編碼邏輯改用固定長度陣列, 減少 GC 與動態分配

## [2.0.3] - 2025-10-16
### Changed
- 更新第三方函式庫

## [2.0.2] - 2025-10-03
### Fixed
- 修正 redmos 執行 Redis 命令時應該無視回應

## [2.0.1] - 2025-10-03
### Added
- 增加非泛型版本的 Raven 訊息工具
### Changed
- 修改泛型版本的 Raven 訊息工具名稱

## [2.0.0] - 2025-10-01
### Deleted
- 移除標籤管理器
- 移除度量管理器
- 移除信號調度管理器
### Changed
- 更新測試工具
- 更新加解密工具
- 更新 IAP 工具
- 更新配置管理器
- 更新執行緒管理器
- 更新日誌管理器
- 更新協助工具
- 更新實體機制
- 更新網路機制
- 更新訊息處理機制
- 更新資料庫機制

## [1.1.57] - 2025-08-13
### Fixed
- 修正 Percent 的計算有問題時不吐 panic, 改回傳0

## [1.1.56] - 2025-08-12
### Fixed
- 修正 Percent 的 Calc 系列函式溢值錯誤

## [1.1.55] - 2025-08-11
### Added
- Entity 新增 ExistMessage 函式

## [1.1.54] - 2025-05-28
### Deleted
- 移除 StringJoin 函式
### Changed
- StringDisplayLength 函式改用第三方函式庫計算

## [1.1.53] - 2025-04-15
### Added
- trials 的 ProtoListMatch 函式, 用來檢查訊息列表中是否有目標訊息

## [1.1.52] - 2025-03-11
### Added
- helps 的 Dice 組件新增 RandOnce 函式
### Changed
- helps 的 Dice 組件擲骰時改用二分搜尋法提升效率

## [1.1.51] - 2025-01-22
### Changed
- 變更 redmos 的 Incr, QPopAll, QPeek 的資料儲存方式, 使其與 Get, Set 等一致

## [1.1.50] - 2025-01-21
### Fixed
- 修正建立專門用於次要資料庫的索引

## [1.1.49] - 2025-01-21
### Added
- redmos 新增 QPush, QPop, QPopAll, QPeek 資料庫函式
- redmos 新增 MinorIndex 函式
- trials 新增 RedisCompareList 函式
### Deleted
- 移除 redmos.Metaer 介面中的 MinorField 函式
### Changed
- 簡化 redmos 次要資料庫儲存程序
- 變更 redmos.Index 結構
### Fixed
- 修正 redmos 次要資料庫批量操作完成後未清除操作列表的問題, 避免下次使用相同物件時重複執行已完成的操作

## [1.1.48] - 2025-01-16
### Added
- 增加 RandSource 來取得隨機來源
### Changed
- 隨機工具改用 golang.org/x/exp/rand 來產生亂數

## [1.1.47] - 2025-01-09
### Added
- 新增 Optsz 組件
### Fixed
- 修正 github workflow

## [1.1.46] - 2025-01-06
### Changed
- 更新第三方函式庫
- 需要 go 1.22 以上

## [1.1.45] - 2024-12-25
### Fixed
- 修正 CalculateDays 的時區計算錯誤

## [1.1.44] - 2024-12-25
### Added
- 新增 CalculateDays, CalculateDaysWithBaseline 函式

## [1.1.43] - 2024-12-23
### Added
- 新增 MapKey, MapValue 函式

## [1.1.42] - 2024-12-23
### Changed
- 把 MapToArray 改為 MapFlatten

## [1.1.41] - 2024-12-11
### Changed
- 重整 Proto 測試工具

## [1.1.40] - 2024-12-11
### Added
- 新增 ProtoTypeExist 函式, 內容為原來的 ProtoContains 函式的功能
### Changed
- ProtoContains 函式的功能改為檢查訊息列表是否有符合訊息

## [1.1.39] - 2024-11-22
### Changed
- 更新 IAP 組件, 讓外部可以獲取交易時間

## [1.1.38] - 2024-10-30
### Added
- 新增 ProtoContains 函式

## [1.1.37] - 2024-10-15
### Added
- 新增 DayHourMax, WeekdayMax 常數

## [1.1.36] - 2024-09-27
### Added
- 新增 Yearly, YearlyPrev, YearlyNext 函式

## [1.1.35] - 2024-09-18
### Added
- 新增 redmos.GetEx 與 redmos.SetEx 命令
### Changed
- 移除 redmos.Get 與 redmos.Set 中的 Expire Time

## [1.1.34] - 2024-09-06
### Added
- 在 Redis Get 的時候加上 Expire Time, 並且每次 Get 都會重置

## [1.1.33] - 2024-09-03
### Changed
- 把 EquateApproxProtoTimestamp 搬移到 helps 軟體包中

## [1.1.32] - 2024-09-03
### Added
- 新增 timestamppb.Timestamp 的比對選項

## [1.1.31] - 2024-09-03
### Changed
- proto 比對改用 go-cmp, 並添加比對選項列表

## [1.1.30] - 2024-08-29
### Changed
- 更新 lint 配置
- 更新 redmos 說明

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
- 封包標頭長度改為 4byte
- 封包大小上限改為 1MB

## [1.1.23] - 2024-07-16
### Deleted
- 移除 ctxs 組件

## [1.1.22] - 2024-07-12
### Added
- 混合資料庫新增取得主要資料庫與次要資料庫

## [1.1.21] - 2024-07-12
### Changed
- 更新 IAP 組件, 讓外部有更多的重試與逾時參數可用

## [1.1.20] - 2024-07-10
### Added
- 新增 ProtoEqual, MapToArray 函式

## [1.1.19] - 2024-06-25
### Added
- 新增 FlagszAND, FlagszOR, FlagszXOR 函式

## [1.1.18] - 2024-06-20
### Changed
- 更新 Fit 組件說明
- StdColor 新增失敗旗標機制
- StdColor 新增 Out, Err 函式

## [1.1.17] - 2024-06-06
### Fixed
- 修正 OneStore 的 IAP 驗證錯誤

## [1.1.16] - 2024-06-05
### Fixed
- 修正 OneStore 的 IAP 驗證錯誤

## [1.1.15] - 2024-06-03
### Added
- 新增 OneStore 的 IAP 組件

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
- 新增 CPU 剖析工具

## [1.1.10] - 2024-05-22
### Changed
- 讓 Raven 的測試系列函式會印出錯誤內容

## [1.1.9] - 2024-05-13
### Added
- 新增 ReflectFieldValue 函式

## [1.1.8] - 2024-05-13
### Added
- 新增 WriteFile 函式

### Changed
- 將 WriteTemplate 改為 TemplateBuild, 並且只負責將新增部分添加到字串尾端

## [1.1.7] - 2024-05-13
### Added
- 新增 WriteTemplate 函式

## [1.1.6] - 2024-05-08
### Changed
- 日誌的 Caller 函式新增可選的紀錄簡單函式名稱功能

## [1.1.5] - 2024-05-06
### Added
- 新增 Before, Beforef, Beforefx 函式
- 新增 After, Afterf, Afterfx 函式
- 新增 Betweenf, Betweenfx 函式
- 新增 Overlapff, Overlapffx 函式

## [1.1.4] - 2024-05-03
### Changed
- 保留 Date 函式的年/月/日參數, 只有時/分/秒/毫秒可以被省略

## [1.1.3] - 2024-05-03
### Fixed
- 修正每月到期計算錯誤

## [1.1.2] - 2024-05-03
### Added
- 新增 JsonString 函式
- 新增 ToProtoAny 函式
- 新增 Timef 函式

### Removed
- 移除 cast 工具
- 移除 Beforef 函式
- 移除 Afterf 函式
- 移除 Betweenf 函式
- 移除 SameDay 函式

### Changed
- ProtoAny 函式改名為 FromProtoAny
- ProtoJson 函式改名為 ProtoString
- Date 函式變得更容易使用

## [1.1.1] - 2024-04-30
### Added
- 新增函式 Beforef, Afterf

## [1.1.0] - 2024-04-26
### Added
- 新增單元測試工具集 trial
- 新增信號調度組件 trigger
- 新增訊息處理組件 procs/raven
- 新增產生錯誤工具

### Changed
- 更新網路組件
- 更新實體組件

## [1.0.0] - 2024-04-15