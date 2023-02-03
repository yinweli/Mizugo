using System.Collections.Generic;

namespace Mizugo
{
    /// <summary>
    /// 事件處理器
    /// </summary>
    public class Eventmgr : IEventmgr
    {
        /// <summary>
        /// 事件處理
        /// </summary>
        /// <param name="eventID">事件編號</param>
        /// <param name="param">事件參數</param>
        public void Process(EventID eventID, object param)
        {
            if (data.TryGetValue(eventID, out var result))
                result?.Invoke(param);
        }

        /// <summary>
        /// 新增事件處理
        /// </summary>
        /// <param name="eventID">事件編號</param>
        /// <param name="onEvent">事件處理函式</param>
        public void Add(EventID eventID, OnTrigger onEvent)
        {
            data[eventID] = onEvent;
        }

        /// <summary>
        /// 刪除事件處理
        /// </summary>
        /// <param name="eventID">事件編號</param>
        public void Del(EventID eventID)
        {
            data.Remove(eventID);
        }

        /// <summary>
        /// 處理列表
        /// </summary>
        private Dictionary<EventID, OnTrigger> data = new Dictionary<EventID, OnTrigger>();
    }
}
