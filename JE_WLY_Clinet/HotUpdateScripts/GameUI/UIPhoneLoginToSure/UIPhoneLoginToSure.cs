using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using JEngine.Core;

public class UIPhoneLoginToSure : UIBase<UIModelPhoneLoginToSure, UIViewPhoneLoginToSure>
{
    public override UILayer GetLayer()
    {
        return UILayer.Normal;
    }
    public override void OnInit()
    {
        GetView().btn_Close.onClick.AddListener(delegate {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Close<UIPhoneLoginToSure>();
        });

        GetView().btn_newCode.onClick.AddListener(delegate {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
        });

        GetView().Input_PhoneCode.onValueChanged.AddListener((call) =>
        {
            if (call.Length == 4)
            { 
            }
        });
    }

    private TimeItem JaTime;
    private int TimeNum = 60; 
    public override void OnEnter(params object[] args)
    {
        GetView().Input_PhoneCode.text = "";
        string phoneStr = PlayerPrefsExtend.GetString("PhoneNum");
        GetView().txt_Num.text = $"发送至{phoneStr.Substring(0, 3)}****{phoneStr.Substring(phoneStr.Length - 4)}";

        TimeNum = 60;

        GetView().root_Time.gameObject.SetActive(true);
        GetView().btn_newCode.gameObject.SetActive(false);

        OnClearTime();
        JaTime = GameHelper.Instance.AddInvoke(new TimeItem(TimeNum, 1, () => {
            GetView().txt_time.text = "重新发送" + TimeNum;
            TimeNum--;
            if (TimeNum <= 0)
            {
                GetView().root_Time.gameObject.SetActive(false);
                GetView().btn_newCode.gameObject.SetActive(true);
                OnClearTime();
            }
        })); 
    }

    private void OnClearTime()
    {
        if (JaTime != null)
        {
            GameHelper.Instance.RemoveInvoke(JaTime);
            JaTime = null;
        }
    }

    public override void OnExit()
    {
        OnClearTime();
    }

     public override void OnDestroy()
    {

    }
}