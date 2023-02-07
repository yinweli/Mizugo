using UnityEngine;
using UnityEngine.UI;

public class CopyToClipboard : MonoBehaviour
{
    public Text targetText;

    public void OnClickCopyToClipboard()
    {
        targetText.text.CopyToClipboard();
    }
}
