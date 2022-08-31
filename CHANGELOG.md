# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Planning]

### 實體-組件機制
這是核心機制, 其他的機制都會依附在這個上面  
* 實體
* 實體管理器
* 組件
* 組件管理器
* 組件間訊息傳遞
    * awake, start, dispose, update
* 執行緒安全

### 網路機制
接聽/連接到遠端, 並產生會話物件來進行封包發送與接收處理  
* 網路統一介面
    * 接聽(go)
    * 連接(go)
        * 只能連接一個服務端
    * 會話(go)
    * 客戶(go+cs)
        * 相當於連接組件與會話組件的組合, 專門給客戶端使用
        * 只能連接一個服務端
* tcp組件
    * 接聽(go)
    * 連接(go)
    * 會話(go)
    * 客戶(go+cs)
* udp組件
* http組件
* 網路機制如何接續到封包機制
* 網路機制自動代碼產生的範圍

### 封包機制
當要發送封包或是接收到封包時, 如何依序進行加密/解密, 序列化/反序列化, 路由/反路由處理程序  
* 封包管線介面
    * 可以自訂的序列化/反序列化機制
    * 可以自訂的加解密機制
    * 可以自訂的路由機制
* proto封包管線組件
* 封包機制如何接續到加解密機制
* 封包機制如何接續到路由機制
* 封包機制自動代碼產生的範圍

### 加解密機制
* ????加解密

### 路由機制
當發送封包時, 如何在封包中附上路由資訊  
當接收封包時, 如何解析路由資訊並轉發到正確目的地去  

### 組網機制
如何組織伺服器之間的連接關係  

### 日誌機制
提供接口讓內部的日誌可以輸出到使用者希望的組件上  

### 資料庫機制
如何初始化資料庫組件, 操作資料庫  
* redis
* mysql
* mongo
* 資料庫機制自動代碼產生的範圍