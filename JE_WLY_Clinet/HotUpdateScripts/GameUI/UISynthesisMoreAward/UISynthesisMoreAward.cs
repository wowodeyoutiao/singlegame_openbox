using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using JEngine.Core;
using JEngine.UI;
using JEngine.Helper;

public class UISynthesisMoreAward : UIBase<UIModelSynthesisMoreAward, UIViewSynthesisMoreAward>
{
    public override UILayer GetLayer()
    {
        return UILayer.Normal;
    }

    public override void OnInit()
    {
        GetView().btn_back.onClick.AddListener(() => {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            ClearJa();
            OnBack();
        });
    }

    private void OnBack()
    {
        UITool.Close<UISynthesisMoreAward>();
    }
    public override void OnEnter(params object[] args)
    {
        AudioManager.Instance.PlaySound(AudioKey.Sound.UI_ComReward);
        var Awrads = args[0] as List<pb.ComposeInfo>;

        for (int i = 0; i < GetView().root_Awards.childCount; i++)
        {
            var go = GetView().root_Awards.GetChild(i);
            go.gameObject.SetActive(i < Awrads.Count);
            if (i < Awrads.Count)
            {
                GameUtil.Instance.SetItemIcon((int)Awrads[i].Uid, JUIExtensions.img(go.transform.Find("icon")), JUIExtensions.img(go.transform),JUIExtensions.img(go.transform.Find("icon/nature")));
                JUIExtensions.txt(go,"txt_num").text = GameUtil.Instance.FormatNumWithUnit(Awrads[i].Num);
            }
        }

        int Count = 6;
        Ja = GameHelper.Instance.AddInvoke(new TimeItem(Count, 1, () => {
            Count--;
            GetView().txt_Time.text = Count + "秒后将自动退出";
            if (Count <=0)
                OnBack();
        })); 
    }

    private TimeItem Ja;
    private void ClearJa()
    {
        if (Ja != null)
        {
            GameHelper.Instance.RemoveInvoke(Ja);
            Ja = null;
        }
    }

    public override void OnExit()
    {
        ClearJa();
    }


    public override void OnDestroy()
    {
        ClearJa();
    }
}