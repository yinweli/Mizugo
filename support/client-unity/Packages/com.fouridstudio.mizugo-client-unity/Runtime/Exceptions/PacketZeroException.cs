using System;

namespace Mizugo
{
    /// <summary>
    /// 封包長度為零異常
    /// </summary>
    public class PacketZeroException : Exception
    {
        public PacketZeroException(string name)
            : base(name + " packet size zero") { }
    }
}
