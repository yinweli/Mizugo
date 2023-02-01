using System;

namespace Mizugo
{
    /// <summary>
    /// 網路結束異常, 這通常代表連接已中斷
    /// 不管是客戶端斷線, 或是伺服器斷線都會造成此異常
    /// </summary>
    public class EofException : Exception
    {
        public EofException()
            : base("eof") { }
    }
}
