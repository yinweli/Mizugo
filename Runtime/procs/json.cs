using System;

namespace Mizugo
{
    /// <summary>
    /// json處理器
    /// </summary>
    public class Json : Procmgr
    {
        public override byte[] Encode(object input)
        {
            var message = input as JsonMsg;

            if (message == null)
                throw new Exception("json encode: invalid message");
        }

        public override object Decode(byte[] input)
        {
            throw new System.NotImplementedException();
        }

        public override bool Process(object message)
        {
            throw new System.NotImplementedException();
        }
    }
}
