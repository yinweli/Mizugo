using System;

namespace Mizugo
{
    /// <summary>
    /// 接收標頭異常
    /// </summary>
    public class RecvHeaderException : Exception
    {
        public RecvHeaderException()
            : base("receive header failed") { }
    }
}
