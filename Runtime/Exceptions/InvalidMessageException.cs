using System;

namespace Mizugo
{
    /// <summary>
    /// 非法訊息異常
    /// </summary>
    public class InvalidMessageException : Exception
    {
        public InvalidMessageException(string name) : base(name + " invalid message") { }
    }
}
