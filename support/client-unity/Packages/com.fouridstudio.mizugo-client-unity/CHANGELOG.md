# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Planning]

## [Unrelease]

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
- 修正連線超時檢測

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