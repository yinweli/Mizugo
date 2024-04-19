using System;

namespace Mizugo
{
    /// <summary>
    /// 封包為空異常
    /// </summary>
    public class PacketNullException : Exception
    {
        public PacketNullException(string name)
            : base(name + " packet null") { }
    }
}
