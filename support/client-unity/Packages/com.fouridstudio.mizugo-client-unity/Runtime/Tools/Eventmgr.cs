using System.Collections.Generic;

namespace Mizugo
{
    /// <summary>
    /// 事件處理器
    /// </summary>
    public class Eventmgr : IEventmgr
    {
        public void Process(EventID eventID, object param)
        {
            if (data.TryGetValue(eventID, out var result))
                result?.Invoke(param);
        }

        public void Add(EventID eventID, OnTrigger onEvent)
        {
            data[eventID] = onEvent;
        }

        public void Del(EventID eventID)
        {
            data.Remove(eventID);
        }

        private Dictionary<EventID, OnTrigger> data = new Dictionary<EventID, OnTrigger>();
    }
}
