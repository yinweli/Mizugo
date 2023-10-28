using System;

namespace Mizugo
{
    /// <summary>
    /// 中斷連線異常
    /// </summary>
    public class DisconnectException : Exception
    {
        public DisconnectException() : base("disconnect") { }
    }
}
