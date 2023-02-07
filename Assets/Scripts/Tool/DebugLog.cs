using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;

public class DebugLog : MonoBehaviour
{
    public CanvasGroup canvasGroup;
    public Transform Content;
    public GameObject objLog;

    public int saveLimit = 100;

    public string output = "";
    public string stack = "";

    public string[] FilterStrs;

    void OnEnable()
    {
        Application.logMessageReceived += HandleLog;
    }

    void OnDisable()
    {
        Application.logMessageReceived -= HandleLog;
    }

    void HandleLog(string logString, string stackTrace, LogType type)
    {
        if(CheckNeedFilter(logString))
            return;

        GameObject obj = Instantiate(objLog, Content);
        Text text = obj.GetComponent<Text>();

        text.text = logString;

        output = logString;
        stack = stackTrace;

        if(Content.childCount > saveLimit)
            ClearTrashLog();            
    }

    bool CheckNeedFilter(string logstr)
    {
        foreach(string filter in FilterStrs)
        {
            if(logstr.Contains(filter))
                return true;
        }

        return false;
    }

    void ClearTrashLog()
    {
        int count = Content.childCount - saveLimit;

        for (int i=1;i< count;i++)
            Destroy(Content.GetChild(i).gameObject);
    }

    public void OnClickOpen()
    {
        canvasGroup.alpha = 1;
        canvasGroup.blocksRaycasts = true;
    }

    public void OnClickClose()
    {
        canvasGroup.alpha = 0;
        canvasGroup.blocksRaycasts = false;
    }

    public void OnClickClear()
    {
        PlayerPrefs.DeleteAll();
    }
}
