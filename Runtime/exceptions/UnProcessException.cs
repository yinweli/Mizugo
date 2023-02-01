using System;

namespace Mizugo
{
    /// <summary>
    /// 訊息未處理異常
    /// </summary>
    public class UnProcessException : Exception
    {
        public UnProcessException(int messageID)
            : base(messageID + " unprocess") { }
    }
}
