using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;

public class CanvasScale : MonoBehaviour
{
    // Start is called before the first frame update
    private void Start()
    {
        var canvasScaler = this.transform.GetComponent<CanvasScaler>();
        var resolution = canvasScaler.referenceResolution;
        var screenRatio = Screen.width / Screen.height;
        if (screenRatio > (resolution.x / resolution.y))
            canvasScaler.matchWidthOrHeight = 1;
        else
            canvasScaler.matchWidthOrHeight = 0;


        //var RectTransform = this.transform.GetComponent<RectTransform>();
        //var info = WeChatWASM.WX.GetSystemInfoSync();
        // if (info == null) return ;
        // float py = (float)info.safeArea.top / (float)info.windowHeight;
        // float by = ((float)info.windowHeight - (float)info.safeArea.bottom) / (float)info.windowHeight;

        //RectTransform.offsetMax = new Vector2(RectTransform.offsetMax.x, m_SafeAreaContainer.offsetMax.y - cs.referenceResolution.y * py);
        //RectTransform.offsetMin = new Vector2(RectTransform.offsetMin.x, m_SafeAreaContainer.offsetMin.y + cs.referenceResolution.y * by);

        //s.referenceResolution = new Vector2(cs.referenceResolution.x, cs.referenceResolution.y * (1.0f + py + (1 - by))); 
    }

    // Update is called once per frame
    void Update()
    {
        
    }
}
