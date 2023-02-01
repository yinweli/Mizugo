using Newtonsoft.Json;
using NUnit.Framework;

namespace Mizugo
{
    internal class TestUtil
    {
        /// <summary>
        /// �ѩ�NUnit���Ѫ�Assert.AreEqual�L�k������O����, ���c����
        /// �ҥH�o�̨ϥΧ⪫���ഫ��json�r��Ӥ��
        /// ���O�p�G���󤺦����X, ���M�i��]�����X���Ǥ��P�y����異��
        /// </summary>
        /// <param name="expected">�w������</param>
        /// <param name="actual">��ڪ���</param>
        public static void AreEqualByJson(object expected, object actual)
        {
            Assert.AreEqual(
                JsonConvert.SerializeObject(expected),
                JsonConvert.SerializeObject(actual)
            );
        }
    }
}
