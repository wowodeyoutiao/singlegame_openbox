using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using JEngine.Event;
using JEngine.Helper;
using HotUpdateScripts.ProtoM;

public class UICom_GVGBattle : UIBase<UIModelCom_GVGBattle, UIViewCom_GVGBattle>
{
    public override UILayer GetLayer()
    {
        return UILayer.NoLayer;
    }
    public override void OnInit()
    {
        JEvent.defaultEvent.Register(this);

        foreach (var item in GetView().rootItems)
        {
            item.txt_Name.text = "";
            item.txt_Num.text = "";
        }
    }

    public override void OnEnter(params object[] args)
    {
    }

    public override void OnExit()
    {
        JEvent.defaultEvent.Unregister(this);
    }
    [Subscriber(ThreadMode.Main)]
    private void OnRefresh(GameEvent.PVP_MeleeKillNum _e)
    {
        _e.killInfos.Sort((a, b) => {return b.Num.CompareTo(a.Num); });
        for (int i = 0; i < GetView().rootItems.Length; i++)
        {
            int QuaIdx = _e.killInfos.Count > i ? (_e.killInfos[i].Uid == RoleM.Instance.Uid ? 2 : 7) : 1;
            GetView().rootItems[i].txt_Name.text = GameUtil.Instance.GetStrByQuality(QuaIdx, (_e.killInfos.Count > i ? _e.killInfos[i].Name : ""));
            GetView().rootItems[i].txt_Num.text = _e.killInfos.Count > i ? "»÷É±Êý:" +_e.killInfos[i].Num : "";
        }
    }

    public override void OnDestroy()
    {
        OnExit();
    }
}