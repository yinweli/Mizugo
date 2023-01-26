namespace Mizugo
{
    /// <summary>
    /// 位址資料
    /// </summary>
    internal struct Addr
    {
        /// <summary>
        /// 位址
        /// </summary>
        public string ip;

        /// <summary>
        /// 埠號
        /// </summary>
        public int port;
    }

    /// <summary>
    /// 日誌介面
    /// </summary>
    internal interface Logger
    {
        /// <summary>
        /// 記錄一般訊息
        /// </summary>
        /// <param name="message">訊息字串</param>
        public void Info(string message);

        /// <summary>
        /// 記錄錯誤訊息
        /// </summary>
        /// <param name="message">訊息字串</param>
        public void Error(string message);
    }
}
