using System;
using System.Collections.Generic;

namespace Mizugo
{
    /// <summary>
    /// 訊息編號, 設置為int32以跟proto的列舉類型統一
    /// </summary>
    using MessageID = Int32;

    /// <summary>
    /// 訊息處理器
    /// </summary>
    public abstract class Procmgr : IProcmgr
    {
        public abstract void Process(object input);

        public void Add(MessageID messageID, OnTrigger onProcess)
        {
            data[messageID] = onProcess;
        }

        public void Del(MessageID messageID)
        {
            data.Remove(messageID);
        }

        public OnTrigger Get(MessageID messageID)
        {
            if (data.TryGetValue(messageID, out var result))
                return result;

            return null;
        }

        private Dictionary<MessageID, OnTrigger> data = new Dictionary<MessageID, OnTrigger>();
    }
}
