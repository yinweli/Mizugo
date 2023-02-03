using Newtonsoft.Json;
using System.Threading;
using UnityEngine;

namespace Mizugo
{
    internal class TestUtil
    {
        /// <summary>
        /// 記錄訊息
        /// </summary>
        /// <param name="message">訊息物件</param>
        public static void Log(object message)
        {
            Debug.Log(message);
        }

        /// <summary>
        /// 共用的暫停函式
        /// </summary>
        public static void Sleep()
        {
            Thread.Sleep(200);
        }

        /// <summary>
        /// 利用json來比對物件, 如果物件內有集合, 仍然可能因為集合順序不同造成比對失敗
        /// </summary>
        /// <param name="expected"></param>
        /// <param name="actual"></param>
        /// <returns></returns>
        public static bool EqualsByJson(object expected, object actual)
        {
            var jsonExpected = JsonConvert.SerializeObject(expected);
            var jsonActual = JsonConvert.SerializeObject(actual);
            return jsonExpected.Equals(jsonActual);
        }
    }
}
