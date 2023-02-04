# Mizugo Client Unity
用於與[mizugo]伺服器連線的Unity客戶端網路組件  

# 系統需求
* [unity]2021.3.9f1以上
* [proto]3以上

# 安裝說明
* 開啟Unity的Package Manager
* 安裝`Google Protobuf`組件
    * 如果專案中已經有該組件, 可以跳過此步驟
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
        * 通常是.Net Framework的版本不符合
        * 或是要修改專案根目錄底下的global.json檔案內容設置
    * 完成後, 編譯完成的檔案會放在`csharp/src/Google.Protobuf/bin/Release`中
    * 依照需求(Unity應該會用`net45`)把該版本的檔案複製到Unity專案中`Assets/Plugins'目錄下
        * 由於各個Unity專案的目錄結構都不太一樣, 因此複製目的地不一定會跟此步驟相同
    * 安裝完成
* 安裝`Newtonsoft Json`組件
    * 如果專案中已經有該組件, 可以跳過此步驟
    * 點擊Package Manager的左上角的`+`號
    * 點擊`Add package by name...`
    * 輸入組件名稱
      ```sh
      com.unity.nuget.newtonsoft-json
      ```
    * 點擊add按鈕
    * 等待安裝完成
* 安裝`Mizugo Client Unity`組件
    * 點擊Package Manager的左上角的`+`號
    * 點擊`Add package from git URL...`
    * 輸入安裝路徑
      ```sh
      https://github.com/yinweli/Mizugo.git#clientunity
      ```
    * 點擊add按鈕
    * 等待安裝完成

![install-client-unity]

# 範例專案
簡單的範例位於`Tests/Runtime/TestTCPClient.cs`, 這是單元測試程式碼  
完整的範例位於[mizugo]專案的`clientunity-sample`分支  

# 專案目錄說明

| 目錄                    | 說明                             |
|:------------------------|:---------------------------------|
| Runtime/Cryptos         | 加密/解密相關組件                |
| Runtime/Events          | 事件組件                         |
| Runtime/Exceptions      | 異常組件                         |
| Runtime/Msgs            | 訊息組件                         |
| Runtime/Nets            | 網路組件                         |
| Runtime/Plugins         | 額外的插件(如protobuf)           |
| Runtime/Procs           | 訊息處理器組件                   |
| Tests/Runtime           | 單元測試                         |
| Tests/Runtime/Msgs      | 單元測試使用的測試訊息           |

[mizugo]: https://github.com/yinweli/sheeter
[proto]: https://github.com/protocolbuffers/protobuf

[install-client-unity]: Documentation/Images/install-client-unity.gif