using System;

namespace Mizugo
{
    /// <summary>
    /// '已經啟動'異常
    /// </summary>
    public class AlreadyStartException : Exception
    {
        public AlreadyStartException(string name)
            : base(name + " already start") { }
    }
}
