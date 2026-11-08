using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;

public class UILoginDes : UIBase<UIModelLoginDes, UIViewLoginDes>
{
    public enum DesType 
    {
        Agree,
        Privacy,
    }
    public override UILayer GetLayer()
    {
        return UILayer.Normal;
    }
     public override void OnInit()
    {
        GetView().btn_Close.onClick.AddListener(()=> {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Close<UILoginDes>();
        });
    }

    public override void OnEnter(params object[] args)
    {
        var selectType = (DesType)args[0];
        GetView().txt_Title.text = selectType == DesType.Agree ? "用户协议" : "隐私政策";
        GetView().obj_Agree.SetActive(selectType == DesType.Agree);
        GetView().obj_Privacy.SetActive(selectType == DesType.Privacy);

        GetView().root_Des.transform.localPosition = Vector3.zero;
    }

    public override void OnExit()
    {
        
    }

     public override void OnDestroy()
    {

    }
}