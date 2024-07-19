namespace Mizugo
{
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
    }

    /// <summary>
    /// 客戶端介面
    /// </summary>
    public interface IClient
    {
        /// <summary>
        /// 連線處理
        /// </summary>
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
        public void Send(object message);

        /// <summary>
        /// 設定事件處理器
        /// </summary>
        public void SetEvent(IEventmgr eventmgr);

        /// <summary>
        /// 新增事件處理
        /// </summary>
        public void AddEvent(EventID eventID, OnTrigger onEvent);

        /// <summary>
        /// 刪除事件處理
        /// </summary>
        public void DelEvent(EventID eventID);

        /// <summary>
        /// 設定訊息處理器
        /// </summary>
        public void SetProc(IProcmgr procmgr);

        /// <summary>
        /// 新增訊息處理
        /// </summary>
        public void AddProcess(int messageID, OnTrigger onProcess);

        /// <summary>
        /// 刪除訊息處理
        /// </summary>
        public void DelProcess(int messageID);

        /// <summary>
        /// 設定編碼/解碼
        /// </summary>
        public void SetCodec(params ICodec[] codec);

        /// <summary>
        /// 設定標頭長度, 連線後就無法變更
        /// </summary>
        public void SetHeaderSize(int size);

        /// <summary>
        /// 設定封包長度, 連線後就無法變更
        /// </summary>
        public void SetPacketSize(int size);

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
        public void Process(EventID eventID, object param);

        /// <summary>
        /// 新增事件處理
        /// </summary>
        public void Add(EventID eventID, OnTrigger onEvent);

        /// <summary>
        /// 刪除事件處理
        /// </summary>
        public void Del(EventID eventID);
    }

    /// <summary>
    /// 訊息處理器介面
    /// </summary>
    public interface IProcmgr
    {
        /// <summary>
        /// 訊息處理
        /// </summary>
        public void Process(object input);

        /// <summary>
        /// 新增訊息處理
        /// </summary>
        public void Add(int messageID, OnTrigger onProcess);

        /// <summary>
        /// 刪除訊息處理
        /// </summary>
        public void Del(int messageID);

        /// <summary>
        /// 取得訊息處理
        /// </summary>
        public OnTrigger Get(int messageID);
    }

    /// <summary>
    /// 編碼/解碼介面
    /// </summary>
    public interface ICodec
    {
        /// <summary>
        /// 編碼處理
        /// </summary>
        public object Encode(object input);

        /// <summary>
        /// 解碼處理
        /// </summary>
        public object Decode(object input);
    }

    /// <summary>
    /// 處理函式類型
    /// </summary>
    public delegate void OnTrigger(object param);
}
