using System;

namespace Mizugo
{
    /// <summary>
    /// 封包長度超過上限異常
    /// </summary>
    public class PacketLimitException : Exception
    {
        public PacketLimitException(string name)
            : base(name + " packet size limit") { }
    }
}
