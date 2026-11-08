using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using HotUpdateScripts.ProtoM;

public class UIPayBindToSure : UIBase<UIModelPayBindToSure, UIViewPayBindToSure>
{
    public override UILayer GetLayer()
    {
        return UILayer.Normal;
    }
    public override void OnInit()
    {
        GetView().btn_Close.onClick.AddListener(delegate {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Close<UIPayBindToSure>();
        });

        GetView().btn_Change.onClick.AddListener(delegate {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Close<UIPayBindToSure>();
            UITool.Open<UIPayBind>();
        });

        GetView().btn_Sure.onClick.AddListener(delegate {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            HttpC.OnExtractAccountUpdate(2, Account, Name, (_r)=>{
                if (_r.code == 0 && !isNullGo)
                {
                    UITool.Close<UIPayBindToSure>();
                }
            });

        });

    }

    private string Account;
    private string Name;

    public override void OnEnter(params object[] args)
    {
        Account = args[0] as string;
        Name = args[1] as string;
        GetView().txt_Info.text = $"姓名:{Name}\n账号:{Account}";
    }

    public override void OnExit()
    {
        
    }

     public override void OnDestroy()
    {

    }
}