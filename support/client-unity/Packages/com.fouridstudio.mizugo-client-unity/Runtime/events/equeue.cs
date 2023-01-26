using System.Collections.Concurrent;

namespace Mizugo
{
    /// <summary>
    /// 事件佇列
    /// </summary>
    internal class EQueue
    {
        /// <summary>
        /// 事件資料
        /// </summary>
        public class Data
        {
            /// <summary>
            /// 事件編號
            /// </summary>
            public EventID eventID;

            /// <summary>
            /// 事件參數
            /// </summary>
            public object param;
        }

        /// <summary>
        /// 新增事件
        /// </summary>
        /// <param name="eventID">事件編號</param>
        /// <param name="param">事件參數</param>
        public void Enqueue(EventID eventID, object param)
        {
            queue.Enqueue(new Data { eventID = eventID, param = param });
        }

        /// <summary>
        /// 取出事件
        /// </summary>
        /// <param name="data">事件物件</param>
        /// <returns>true表示成功, false則否</returns>
        public bool Dequeue(out Data data)
        {
            return queue.TryDequeue(out data);
        }

        /// <summary>
        /// 事件列表
        /// </summary>
        private ConcurrentQueue<Data> queue = new ConcurrentQueue<Data>();
    }
}
