# Mizugo Client Unity
用於與[mizugo]伺服器連線的[unity]客戶端網路組件  

# 系統需求
- [unity]2021.3.9f1以上
- [proto]3以上

# 安裝說明
- 安裝`Google Protobuf`組件
    - 如果你的專案已經包含此組件, 可以跳過此步驟
    - 方法一: 自行編譯
        - 確保已安裝Visual Studio Community, 可從下面的位址下載
          ```sh
          https://visualstudio.microsoft.com/zh-hant/vs/community/
          ```
        - 訪問`Google Protobuf`的GitHub頁面
          ```sh
          https://github.com/protocolbuffers/protobuf
          ```
        - 從`Release`部分選擇所需的版本, 目前使用版本是`Protocol Buffers v21.12`, 下載`protobuf-all-21.12.zip`後解壓縮檔案
        - 打開`csharp/src/Google.Protobuf.sln`, 在Visual Studio Community中選擇`建置/批次建置`, 勾選`Google.Protobuf`的Release並建置專案
        - 若遇到任何問題, 需自行排除; 若有如下錯誤信息
          ```sh
          Detailed Information: Unable to locate the .NET Core SDK. Check that it is installed and that the version specified in global.json (if any) matches the installed version.
          ```
          這可能是因為SDK版本不匹配, 通過命令行執行`dotnet --version`檢查並適當調整global.json文件中的版本號
        - 編譯完成的檔案將位於`csharp/src/Google.Protobuf/bin/Release`, 根據[unity]需求(通常是`net45`), 將相應文件複製到`Assets/Plugins`目錄下
    - 方法二: 從Package Manager安裝
        - 在[unity]的Package Manager中點擊左上角的 + 號
        - 選擇`Add package from git URL...`, 輸入以下路徑
          ```sh
          https://github.com/yinweli/Mizugo.git#proto-unity
          ```
        - 點擊add按鈕, 等待安裝完成
- 安裝`Newtonsoft Json`組件
    - 如果你的專案已經包含此組件, 可以跳過此步驟
    - 在[unity]的Package Manager中點擊左上角的 + 號
    - 選擇`Add package by name...`, 輸入組件名稱
      ```sh
      com.unity.nuget.newtonsoft-json
      ```
    - 點擊add按鈕, 等待安裝完成
- 安裝`Mizugo Client Unity`組件
    - 在[unity]的Package Manager中點擊左上角的 + 號
    - 選擇`Add package from git URL...`, 輸入安裝路徑
      ```sh
      https://github.com/yinweli/Mizugo.git#client-unity
      ```
    - 點擊add按鈕, 等待安裝完成

![install-client-unity]

# 範例專案
位於`Tests/Runtime/TestTCPClient.cs`的單元測試程式碼示範基本使用, 完整範例位於[mizugo]專案的`client-unity-sample`分支

# 專案說明
客戶端組件由四個核心組件組成: [網路組件](#網路組件), [事件組件](#事件組件), [訊息處理組件](#訊息處理組件), [編碼/解碼流程](#編碼/解碼流程)

## 網路組件
負責連線、斷線、更新處理、訊息傳送, 事件處理, 編碼/解碼流程, 需要定時執行Update函式來處理事件與訊息  
目前有下列網路組件可選  
- TCPClient, 使用TCP來連線到伺服器
  ```cs
  var client = new TCPClient();
  var eventmgr = new Eventmgr();
  var process = new ProcJson();
  
  client.SetEvent(eventmgr);
  client.SetProc(process);
  client.SetCodec(process, new DesCBC(PaddingMode.PKCS7, key, iv), new Base64()); // 設定編碼/解碼流程, 這裡設定了依序做ProcJson, desCBC, base64的編碼/解碼
  
  client.AddEvent(...);   // 註冊事件處理
  client.AddProcess(...); // 註冊訊息處理
  client.Connect(...);    // 進行連線
  client.Update();        // 更新處理
  client.Disconnect();    // 進行斷線
  ```

## 事件組件
負責處理網路事件, 避免使用者直接處理多執行緒問題  
目前有下列事件組件可選  
- Eventmgr, 標準的事件處理器

## 訊息處理組件
負責訊息的實際處理, 不同的組件使用的封包結構不同, 因此訊息處理函式無法混用  
目前有下列訊息處理組件可選  
- ProcJson, 使用JsonMsg結構通訊, 訊息定義在 support/proto/mizugo/msg-go/msgs-json/jsonmsg.go, support/proto/mizugo/msg-cs/msgs-json/Jsonmsg.cs
- ProcProto, 使用ProtoMsg結構通訊, 訊息定義在 support/proto/mizugo/protomsg.proto

## 編碼/解碼流程
定義訊息傳送和接收的編碼及解碼過程, 應與伺服器的設定一致

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

| 異常類型              | 異常名稱             | 說明                                                                 |
|:----------------------|:---------------------|:---------------------------------------------------------------------|
| AlreadyStartException | '已經啟動'異常       | 發生在對已經完成連線或是正在連線中的客戶端組件執行連線函式時         |
| DisconnectException   | 中斷連線異常         | 發生在斷線時                                                         |
| PacketLimitException  | 封包長度超過上限異常 | 發生在接收/傳送過大封包時                                            |
| PacketNullException   | 封包為空異常         | 發生在接收/傳送流程中遭遇空指標時                                    |
| PacketZeroException   | 封包長度為零異常     | 發生在接收/傳送空封包時                                              |
| RecvHeaderException   | 接收標頭異常         | 發生在接收標頭但是與預期長度不一致時                                 |
| RecvPacketException   | 接收封包異常         | 發生在接收封包但是與預期長度不一致時                                 |
| UnprocessException    | 訊息未處理異常       | 發生在接收訊息但沒有對應的訊息處理函式時, 這個錯誤看情況是可以忽略的 |

除了UnprocessException以外, 其他的異常都是嚴重的錯誤, 應該要中斷連線  
另外執行中如果網路失敗, 會收到內建的SocketException異常, 可以通過把異常轉型為SocketException並且判讀其中的SocketErrorCode來決定處理方式  
SocketErrorCode是個列舉, 內容可以參考[socket-error-enum]或是[socket-error-code]  

# 專案目錄說明
| 目錄               | 說明                   |
|:-------------------|:-----------------------|
| Runtime/Clients    | 網路相關組件           |
| Runtime/Cryptos    | 加密/解密相關組件      |
| Runtime/Exceptions | 異常組件               |
| Runtime/Msgs       | 訊息組件               |
| Tests/Runtime      | 單元測試               |
| Tests/Runtime/Msgs | 單元測試使用的測試訊息 |

[mizugo]: https://github.com/yinweli/mizugo
[proto]: https://github.com/protocolbuffers/protobuf
[unity]: https://unity.com/
[socket-error-enum]: https://learn.microsoft.com/zh-tw/dotnet/api/system.net.sockets.socketerror?view=netframework-4.8
[socket-error-code]: https://learn.microsoft.com/zh-tw/windows/win32/winsock/windows-sockets-error-codes-2

[install-client-unity]: Documentation/Images/install-client-unity.gif