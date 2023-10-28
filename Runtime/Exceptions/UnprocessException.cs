using System;

namespace Mizugo
{
    /// <summary>
    /// 訊息未處理異常
    /// </summary>
    public class UnprocessException : Exception
    {
        public UnprocessException(params int[] messageID) : base(string.Join(",", messageID) + " unprocess") { }
    }
}
