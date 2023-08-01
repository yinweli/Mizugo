using System;

namespace Mizugo
{
    /// <summary>
    /// 接收封包異常
    /// </summary>
    public class RecvPacketException : Exception
    {
        public RecvPacketException()
            : base("receive packet failed") { }
    }
}
