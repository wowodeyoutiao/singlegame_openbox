using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;

public class UIPhoneLogin : UIBase<UIModelPhoneLogin, UIViewPhoneLogin>
{
    public override UILayer GetLayer()
    {
        return UILayer.Normal;
    }
    public override void OnInit()
    {
        GetView().btn_Close.onClick.AddListener(delegate {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Close<UIPhoneLogin>();
        });

        GetView().btn_agree.onClick.AddListener(delegate {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Open<UILoginDes>(null, UILoginDes.DesType.Agree);
        });

        GetView().btn_yinsi.onClick.AddListener(delegate {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Open<UILoginDes>(null, UILoginDes.DesType.Privacy);
        });

        GetView().btn_Sure.onClick.AddListener(delegate {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            if (PlayerPrefsExtend.GetInt("Agree") == 0)
            {
                UITool.ShowTips("请先阅读并同意《用户协议》、《隐私政策》");
                return;
            }
            if (string.IsNullOrEmpty(GetView().Input_Phone.text) || GetView().Input_Phone.text.Length != 11)
            {
                UITool.ShowTips("手机号错误");
                return;
            }
            PlayerPrefsExtend.SetString("PhoneNum", GetView().Input_Phone.text);
            UITool.Open<UIPhoneLoginToSure>();
        });

        GetView().tog_agree.onValueChanged.AddListener((isOn) => {
            PlayerPrefsExtend.SetInt("Agree", isOn ? 1 : 0);
        });
    }

    public override void OnEnter(params object[] args)
    {
        GetView().tog_agree.isOn = PlayerPrefsExtend.GetInt("Agree") == 0;
        GetView().Input_Phone.text = PlayerPrefsExtend.GetString("PhoneNum");
    }

    public override void OnExit()
    {
        
    }

     public override void OnDestroy()
    {

    }
}