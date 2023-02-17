/// <summary>
/// 要求Json
/// </summary>
public class MJsonQ
{
    /// <summary>
    /// 傳送時間
    /// </summary>
    public long Time;
}

/// <summary>
/// 回應Json
/// </summary>
public class MJsonA
{
    /// <summary>
    /// 來源訊息
    /// </summary>
    public MJsonQ From;

    /// <summary>
    /// 錯誤編號
    /// </summary>
    public int ErrID;

    /// <summary>
    /// 封包計數
    /// </summary>
    public long Count;
}
