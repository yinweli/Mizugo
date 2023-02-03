using System;

namespace Mizugo
{
    /// <summary>
    /// 接收封包異常
    /// </summary>
    public class ReceiveException : Exception
    {
        public ReceiveException()
            : base("receive failed") { }
    }
}
