using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;

public class UIPayBind : UIBase<UIModelPayBind, UIViewPayBind>
{
    public override UILayer GetLayer()
    {
        return UILayer.Normal;
    }
    public override void OnInit()
    {
        GetView().btn_Close.onClick.AddListener(delegate {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Close<UIPayBind>();
        });

        GetView().btn_Sure.onClick.AddListener(delegate {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            if (string.IsNullOrEmpty(GetView().Input_Acc.text))
            {
                UITool.ShowTips("账号不能为空");
                return;
            }

            if (string.IsNullOrEmpty(GetView().Input_Name.text))
            {
                UITool.ShowTips("姓名不能为空");
                return;
            }
            UITool.Close<UIPayBind>();
            UITool.Open<UIPayBindToSure>(null, GetView().Input_Acc.text, GetView().Input_Name.text);
        });
    }

    public override void OnEnter(params object[] args)
    {
        GetView().Input_Name.text = "";
        GetView().Input_Acc.text = "";
    }

    public override void OnExit()
    {
        
    }

     public override void OnDestroy()
    {

    }
}