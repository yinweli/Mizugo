# Mizugo Client Unity
用於與[mizugo]伺服器連線的[unity]客戶端網路組件  

# 系統需求
* [unity]2021.3.9f1以上
* [proto]3以上

# 安裝說明
* 安裝`Google Protobuf`組件
    * 如果專案中已經有該組件, 可以跳過此步驟
    * 第一種方法: 自行編譯
        * 此步驟需要事先安裝Visual Studio Community, 可以到以下位址下載
          ```sh
          https://visualstudio.microsoft.com/zh-hant/vs/community/
          ```
        * 使用瀏覽器到達`Google Protobuf`的Github主頁
          ```sh
          https://github.com/protocolbuffers/protobuf
          ```
        * 從右邊的`Release`找到所需的版本, 目前使用`Protocol Buffers v21.12`
        * 從Assets列表中下載`protobuf-all-xx.xx.zip`, 目前使用`protobuf-all-21.12.zip`
        * 解壓縮檔案
        * 開啟protobuf專案目錄中的`csharp/src/Google.Protobuf.sln`檔案
        * 選擇功能列中的`建置/批次建置`, 並把`Google.Protobuf`的Release打勾, 並建置專案
            * 如果建置過程中有發生問題, 得要自己排除問題了QQ
            * 如果有輸出如下錯誤
              ```sh
              Detailed Information: Unable to locate the .NET Core SDK. Check that it is installed and that the version specified in global.json (if any) matches the installed version.
              ```
              這是因為Visual Studio Community使用的dotnet core sdk版本與protobuf專案裡使用的sdk版本不一致導致的
              可通過在命令行中執行
              ```sh
              dotnet --version
              ```
              來查看當前系統使用的dotnet core sdk版本, 然後以此修改protobuf專案根目錄底下的global.json檔案中的sdk版號來解決
        * 完成後, 編譯完成的檔案會放在`csharp/src/Google.Protobuf/bin/Release`中
        * 依照需求([unity]應該會用`net45`)把該版本的檔案複製到[unity]專案中`Assets/Plugins'目錄下
            * 由於各個[unity]專案的目錄結構都不太一樣, 因此複製目的地不一定會跟此步驟相同
        * 安裝完成
    * 第二種方法: 從Package Manager安裝
        * 開啟[unity]的Package Manager
        * 點擊Package Manager的左上角的`+`號
        * 點擊`Add package from git URL...`
        * 輸入安裝路徑
          ```sh
          https://github.com/yinweli/Mizugo.git#proto-unity
          ```
        * 點擊add按鈕
        * 等待安裝完成
* 安裝`Newtonsoft Json`組件
    * 如果專案中已經有該組件, 可以跳過此步驟
    * 開啟[unity]的Package Manager
    * 點擊Package Manager的左上角的`+`號
    * 點擊`Add package by name...`
    * 輸入組件名稱
      ```sh
      com.unity.nuget.newtonsoft-json
      ```
    * 點擊add按鈕
    * 等待安裝完成
* 安裝`Mizugo Client Unity`組件
    * 開啟[unity]的Package Manager
    * 點擊Package Manager的左上角的`+`號
    * 點擊`Add package from git URL...`
    * 輸入安裝路徑
      ```sh
      https://github.com/yinweli/Mizugo.git#client-unity
      ```
    * 點擊add按鈕
    * 等待安裝完成

![install-client-unity]

# 範例專案
簡單的範例位於`Tests/Runtime/TestTCPClient.cs`, 這是單元測試程式碼  
完整的範例位於[mizugo]專案的`client-unity-sample`分支  

# 專案說明
客戶端組件由三個核心組件組成: [網路組件](#網路組件), [事件組件](#事件組件), [訊息處理組件](#訊息處理組件)  

## 網路組件
網路組件繼承了IClient介面, 是最核心的組件  
負責連線處理, 斷線處理, 更新處理, 傳送訊息, 新增/刪除事件處理, 新增/刪除訊息處理  
建立網路組件時需要指定使用哪個事件組件與訊息處理組件  
連線後需要定時執行Update函式來執行事件與訊息處理  
目前有下列網路組件可選  
* TCPClient: 
    - 使用TCP來連線到伺服器
    - 建立範例
      ```cs
      var client = new TCPClient(new Eventmgr(), new JsonProc());
      ```

## 事件組件
事件組件繼承了IEventmgr介面, 負責網路組件或是訊息處理組件與外部的溝通, 避免使用者直接面對多執行緒  
目前有下列事件組件可選  
* Eventmgr
    - 標準的事件處理器

## 訊息處理組件
訊息處理組件繼承了IProcmgr介面, 負責封包編碼/解碼, 管理訊息處理函式  
不同的訊息處理組件使用的封包結構不同, 因此訊息處理函式也無法通用  
目前有下列訊息處理組件可選  
* JsonProc
    - 沒有使用加密技術, 所以安全性很低, 僅用於傳送簡單訊息或是傳送密鑰使用
    - 封包結構: support/proto/mizugo/msg-cs/msgs-json/Jsonmsg.cs
    - 編碼方式: json編碼成位元陣列, 再通過base64編碼
    - 解碼方式: base64解碼, 再通過json解碼
* ProtoProc
    - 沒有使用加密技術, 所以安全性很低, 僅用於傳送簡單訊息或是傳送密鑰使用
    - 封包結構: support/proto/mizugo/protomsg.proto
    - 編碼方式: protobuf編碼成位元陣列, 再通過base64編碼
    - 解碼方式: base64解碼, 再通過protobuf解碼
* PListProc
    - 採用des-cbc加密, 安全性較高, 適合用來傳送一般封包, 使用時需要設定密鑰以及初始向量
    - 採用複數訊息設計, 因此封包內可以填入多個訊息來跟伺服器溝通
    - 封包結構: support/proto/mizugo/plistmsg.proto
    - 編碼方式: protobuf編碼成位元陣列, 再通過des加密
    - 解碼方式: des解密, 再通過protobuf解碼

## 封包限制
[mizugo]的封包長度最大不能超過65535個位元組, 相當於64K位元資料  

## 事件
客戶端組件內發生的事件都定義在EventID列舉中  

| 名稱       | 說明                                           |
|:-----------|:-----------------------------------------------|
| Connect    | 連線事件, 參數是null                           |
| Disconnect | 斷線事件, 參數是null                           |
| Recv       | 接收事件, 當接收並處理完封包後執行, 參數是null |
| Send       | 傳送事件, 當傳送封包完畢後執行, 參數是null     |
| Error      | 錯誤事件, 參數是Exception                      |

## 異常
客戶端組件內發生的異常除了.Net原有異常之外, 還有下列新增的異常  

| 異常類型                | 異常名稱             | 說明                                                                 |
|:------------------------|:---------------------|:---------------------------------------------------------------------|
| AlreadyStartException   | '已經啟動'異常       | 發生在對已經完成連線或是正在連線中的客戶端組件執行連線函式時         |
| InvalidMessageException | 非法訊息異常         | 發生在客戶端與伺服器雙方使用的訊息處理器不同時                       |
| PacketLimitException    | 封包長度超過上限異常 | 發生在接收/傳送過大封包時                                            |
| PacketZeroException     | 封包長度為零異常     | 發生在接收/傳送空封包時                                              |
| RecvHeaderException     | 接收標頭異常         | 發生在接收標頭但是與預期長度不一致時                                 |
| RecvPacketException     | 接收封包異常         | 發生在接收封包但是與預期長度不一致時                                 |
| UnprocessException      | 訊息未處理異常       | 發生在接收訊息但沒有對應的訊息處理函式時, 這個錯誤看情況是可以忽略的 |

除了UnprocessException以外, 其他的異常都是嚴重的錯誤, 應該要中斷連線  

# 專案目錄說明

| 目錄                    | 說明                             |
|:------------------------|:---------------------------------|
| Runtime/Cryptos         | 加密/解密相關組件                |
| Runtime/Events          | 事件組件                         |
| Runtime/Exceptions      | 異常組件                         |
| Runtime/Msgs            | 訊息組件                         |
| Runtime/Nets            | 網路組件                         |
| Runtime/Procs           | 訊息處理器組件                   |
| Tests/Runtime           | 單元測試                         |
| Tests/Runtime/Msgs      | 單元測試使用的測試訊息           |

[mizugo]: https://github.com/yinweli/mizugo
[proto]: https://github.com/protocolbuffers/protobuf
[unity]: https://unity.com/

[install-client-unity]: Documentation/Images/install-client-unity.gif