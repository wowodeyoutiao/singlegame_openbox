using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using System;

public class UIMessageDestoryRole : UIBase<UIModelMessageDestoryRole, UIViewMessageDestoryRole>
{
    public override UILayer GetLayer()
    {
        return UILayer.Normal;
    }
    public override void OnInit()
    {
        GetView().Btn_No.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Close<UIMessageDestoryRole>();
        });

        GetView().Btn_Yes.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            if (GetView().Intput.text != "删除")
            {
                UITool.ShowTips("请正常输入文字内容！");
                return;
            }
            UITool.Close<UIMessageDestoryRole>();
            CallBack?.Invoke();
        });
    }

    private Action CallBack;
    public override void OnEnter(params object[] args)
    {
        GetView().Intput.text = "";
        CallBack = args[0] as Action;
    }

    public override void OnExit()
    {
        
    }

     public override void OnDestroy()
    {

    }
}