using Newtonsoft.Json;
using NUnit.Framework;

namespace Mizugo
{
    internal class TestUtil
    {
        /// <summary>
        /// 由於NUnit提供的Assert.AreEqual無法比對類別物件, 結構物件
        /// 所以這裡使用把物件轉換為json字串來比對
        /// 但是如果物件內有集合, 仍然可能因為集合順序不同造成比對失敗
        /// </summary>
        /// <param name="expected">預期物件</param>
        /// <param name="actual">實際物件</param>
        public static void AreEqualByJson(object expected, object actual)
        {
            Assert.AreEqual(
                JsonConvert.SerializeObject(expected),
                JsonConvert.SerializeObject(actual)
            );
        }
    }
}
