using System;
using System.Diagnostics;
using System.Security.Cryptography;
using Mizugo;
using UnityEngine;

/// <summary>
/// 客戶端組件範例, 使用Proto訊息處理器, 編碼/解碼流程採用ProcProto, DesCBC以及Base64
/// 程式會在Awake時初始化內部組件, 在Start時連線到伺服器, 在Update時更新客戶端組件
/// 連線成功後, 在OnConnect時傳送MProtoQ訊息到伺服器, 等待伺服器的回應
/// 當伺服器回應MProtoA訊息時, 在ProcMProtoA處理它並顯示訊息, 訊息顯示完畢後就斷線
/// 此範例需要配合Mizugo專案的測試伺服器才能正常運作
/// </summary>
public class SampleProto : MonoBehaviour
{
    private void Awake()
    {
        var eventmgr = new Eventmgr();
        var process = new ProcProto();

        client = new TCPClient();
        client.SetEvent(eventmgr);
        client.SetProc(process);
        client.SetCodec(process, new DesCBC(PaddingMode.PKCS7, key, key), new Base64());
        client.AddEvent(EventID.Connect, OnConnect);
        client.AddEvent(EventID.Disconnect, OnDisconnect);
        client.AddEvent(EventID.Recv, OnRecv);
        client.AddEvent(EventID.Send, OnSend);
        client.AddEvent(EventID.Error, OnError);
        client.AddProcess((int)MsgID.ProtoA, ProcMProtoA);
        stopwatch = new Stopwatch();
    }

    private void Start()
    {
        try
        {
            client?.Connect(host, port);
        } // try
        catch (Exception e)
        {
            Log("connect to " + host + ":" + port + " failed: " + e);
        } // catch
    }

    private void Update()
    {
        client?.Update();
    }

    /// <summary>
    /// 連線通知, param總是為null
    /// </summary>
    private void OnConnect(object param)
    {
        Log("connect to " + host + ":" + port + " success");
        stopwatch?.Start();
        SendMProtoQ();
    }

    /// <summary>
    /// 斷線通知, param總是為null
    /// </summary>
    private void OnDisconnect(object param)
    {
        Log("disconnect");
    }

    /// <summary>
    /// 接收封包通知, param總是為null
    /// </summary>
    private void OnRecv(object param)
    {
        Log("recv packet");
    }

    /// <summary>
    /// 傳送封包通知, param總是為null
    /// </summary>
    private void OnSend(object param)
    {
        Log("send packet");
    }

    /// <summary>
    /// 錯誤通知, param是Exception物件
    /// </summary>
    private void OnError(object param)
    {
        Log(param);
    }

    /// <summary>
    /// 訊息處理: MProtoA
    /// 處理訊息時, 首要就是要把param物件轉換為訊息結構
    /// 當使用Proto訊息處理器時, 可以通過ProcProto.Unmarshal函式來幫助轉換為訊息結構
    /// 由於一個訊息處理函式只針對一個訊息處理, 因此可以確定要轉換的訊息結構類型
    /// 如果ProcProto.Unmarshal或是訊息處理函式拋出異常, 會由客戶端組件負責捕獲
    /// </summary>
    private void ProcMProtoA(object param)
    {
        ProcProto.Unmarshal<MProtoA>(param, out var messageID, out var message);
        var duration = stopwatch.ElapsedMilliseconds - message.From.Time;
        var count = message.Count;
        var errID = message.ErrID;

        Log(">>> duration: " + duration + ", count: " + count + ", errID: " + errID);
        client.Disconnect();
    }

    /// <summary>
    /// 訊息傳送: MProtoQ
    /// 傳送訊息時, 首要就是要建立訊息結構並填寫好各個欄位
    /// 當使用Proto訊息處理器時, 可以通過ProcProto.Marshal函式來幫助建立訊息結構
    /// ProcProto.Marshal函式需要填寫訊息編號以及訊息物件來建立訊息結構
    /// 這個範例函式專門用來傳送MProtoQ訊息, 訊息中只有一個Time欄位用來填寫當前時間的毫秒
    /// 最後用客戶端組件把訊息傳送出去
    /// </summary>
    private void SendMProtoQ()
    {
        var message = ProcProto.Marshal((int)MsgID.ProtoQ, new MProtoQ { Time = stopwatch.ElapsedMilliseconds });

        client.Send(message);
    }

    /// <summary>
    /// 輸出日誌
    /// </summary>
    private void Log(object message)
    {
        UnityEngine.Debug.Log("sample proto: " + message);
    }

    /// <summary>
    /// 伺服器位址
    /// </summary>
    [SerializeField]
    private string host = string.Empty;

    /// <summary>
    /// 伺服器埠號
    /// </summary>
    [SerializeField]
    private int port = 9002;

    /// <summary>
    /// 密鑰
    /// </summary>
    [SerializeField]
    private string key = "key-####";

    /// <summary>
    /// 客戶端組件
    /// </summary>
    private TCPClient client = null;

    /// <summary>
    /// 計時器
    /// </summary>
    private Stopwatch stopwatch = null;
}
