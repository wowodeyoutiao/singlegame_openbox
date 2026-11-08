using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;

public class UINewFunctionOpen : UIBase<UIModelNewFunctionOpen, UIViewNewFunctionOpen>
{
    public override UILayer GetLayer()
    {
        return UILayer.Base;
    }
    public override void OnInit()
    {
        GetView().btn_Close.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Close<UINewFunctionOpen>();
        });
    }

    string content;
    public override void OnEnter(params object[] args)
    {
        content = (string)args[0];
        GetView().txt_Desc.text = content;
        GetView().scrollRect.verticalNormalizedPosition = 1f;
    }

    public override void OnExit()
    {
        
    }

    public override void OnDestroy()
    {

    }
}