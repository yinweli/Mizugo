using System.Collections.Generic;

namespace Mizugo
{
    /// <summary>
    /// 訊息處理器
    /// </summary>
    public abstract class Procmgr : IProcmgr
    {
        public abstract void Process(object input);

        public void Add(int messageID, OnTrigger onProcess)
        {
            data[messageID] = onProcess;
        }

        public void Del(int messageID)
        {
            data.Remove(messageID);
        }

        public OnTrigger Get(int messageID)
        {
            if (data.TryGetValue(messageID, out var result))
                return result;

            return null;
        }

        private Dictionary<int, OnTrigger> data = new Dictionary<int, OnTrigger>();
    }
}
