using System;

namespace Mizugo
{
    /// <summary>
    /// 訊息編號, 設置為int32以跟proto的列舉類型統一
    /// </summary>
    using MessageID = Int32;

    /// <summary>
    /// 事件編號
    /// </summary>
    public enum EventID
    {
        /// <summary>
        /// 無效事件
        /// </summary>
        Unknow = 0,

        /// <summary>
        /// 連線事件, 參數是null
        /// </summary>
        Connect,

        /// <summary>
        /// 斷線事件, 參數是null
        /// </summary>
        Disconnect,

        /// <summary>
        /// 接收事件, 參數是null
        /// </summary>
        Recv,

        /// <summary>
        /// 傳送事件, 參數是null
        /// </summary>
        Send,

        /// <summary>
        /// 錯誤事件, 參數是Exception
        /// </summary>
        Error,

        /// <summary>
        /// 訊息事件, 參數是訊息物件; 外部無法使用
        /// </summary>
        Message,
    }

    /// <summary>
    /// 處理函式類型
    /// </summary>
    /// <param name="param">參數物件</param>
    public delegate void OnTrigger(object param);

    /// <summary>
    /// 網路定義
    /// </summary>
    public class Define
    {
        /// <summary>
        /// 標頭長度
        /// </summary>
        public const int headerSize = 2;

        /// <summary>
        /// 封包長度
        /// </summary>
        public const int packetSize = ushort.MaxValue;

        /// <summary>
        /// 緩衝區長度
        /// </summary>
        public const int bufferSize = (headerSize + packetSize) * 10;
    }

    /// <summary>
    /// 客戶端介面
    /// </summary>
    public interface IClient
    {
        /// <summary>
        /// 連線處理
        /// </summary>
        /// <param name="host">連線位址</param>
        /// <param name="port">連線埠號</param>
        public void Connect(string host, int port);

        /// <summary>
        /// 斷線處理
        /// </summary>
        public void Disconnect();

        /// <summary>
        /// 更新處理
        /// </summary>
        public void Update();

        /// <summary>
        /// 傳送訊息
        /// </summary>
        /// <param name="message">訊息物件</param>
        public void Send(object message);

        /// <summary>
        /// 新增事件處理
        /// </summary>
        /// <param name="eventID">事件編號</param>
        /// <param name="onEvent">事件處理函式</param>
        public void AddEvent(EventID eventID, OnTrigger onEvent);

        /// <summary>
        /// 刪除事件處理
        /// </summary>
        /// <param name="eventID">事件編號</param>
        public void DelEvent(EventID eventID);

        /// <summary>
        /// 新增訊息處理
        /// </summary>
        /// <param name="messageID">訊息編號</param>
        /// <param name="onProcess">訊息處理函式</param>
        public void AddProcess(MessageID messageID, OnTrigger onProcess);

        /// <summary>
        /// 刪除訊息處理
        /// </summary>
        /// <param name="messageID">訊息編號</param>
        public void DelProcess(MessageID messageID);

        /// <summary>
        /// 取得連線位址
        /// </summary>
        public string Host { get; }

        /// <summary>
        /// 取得連線埠號
        /// </summary>
        public int Port { get; }

        /// <summary>
        /// 取得是否連線
        /// </summary>
        public bool IsConnect { get; }

        /// <summary>
        /// 取得是否需要處理事件, 連線中或是事件佇列中有事件都會需要處理事件
        /// </summary>
        public bool IsUpdate { get; }
    }

    /// <summary>
    /// 事件處理器介面
    /// </summary>
    public interface IEventmgr
    {
        /// <summary>
        /// 事件處理
        /// </summary>
        /// <param name="eventID">事件編號</param>
        /// <param name="param">事件參數</param>
        public void Process(EventID eventID, object param);

        /// <summary>
        /// 新增事件處理
        /// </summary>
        /// <param name="eventID">事件編號</param>
        /// <param name="onEvent">事件處理函式</param>
        public void Add(EventID eventID, OnTrigger onEvent);

        /// <summary>
        /// 刪除事件處理
        /// </summary>
        /// <param name="eventID">事件編號</param>
        public void Del(EventID eventID);
    }

    /// <summary>
    /// 訊息處理器介面
    /// </summary>
    public interface IProcmgr
    {
        /// <summary>
        /// 封包編碼
        /// </summary>
        /// <param name="input">輸入物件</param>
        /// <returns>結果物件</returns>
        public byte[] Encode(object input);

        /// <summary>
        /// 封包解碼
        /// </summary>
        /// <param name="input">輸入物件</param>
        /// <returns>結果物件</returns>
        public object Decode(byte[] input);

        /// <summary>
        /// 訊息處理
        /// </summary>
        /// <param name="input">輸入物件</param>
        public void Process(object input);

        /// <summary>
        /// 新增訊息處理
        /// </summary>
        /// <param name="messageID">訊息編號</param>
        /// <param name="onProcess">訊息處理函式</param>
        public void Add(MessageID messageID, OnTrigger onProcess);

        /// <summary>
        /// 刪除訊息處理
        /// </summary>
        /// <param name="messageID">訊息編號</param>
        public void Del(MessageID messageID);

        /// <summary>
        /// 取得訊息處理
        /// </summary>
        /// <param name="messageID">訊息編號</param>
        /// <returns>訊息處理函式</returns>
        public OnTrigger Get(MessageID messageID);
    }
}
