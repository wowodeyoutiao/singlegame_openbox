using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using JEngine.Core;
using HotUpdateScripts.ProtoM;
using HotUpdateScripts.Game;

public class UIPVPBeginTime : UIBase<UIModelPVPBeginTime, UIViewPVPBeginTime>
{
    public enum ShowType
    {
        Begin,
        DieToLife,
        MyTeamFirstAttack,
        OtherTeamFirstAttack,
    }

    public override UILayer GetLayer()
    {
        return UILayer.Normal;
    }

    public override void OnInit()
    {

    }

    public override void OnEnter(params object[] args)
    {
        var sType = (ShowType)args[0];
        int Count = sType == ShowType.Begin ? 3 : CFG_XMLConfig.xmlConfig.GVGInfo.alivetime;
        string timeStr = "";
        switch (sType)
        {
            case ShowType.Begin:
                Count = 3;
                GetView().txt_tip.text = "";
                timeStr = "准备战斗:";
                break;
            case ShowType.DieToLife:
                Count = CFG_XMLConfig.xmlConfig.GVGInfo.alivetime;
                GetView().txt_tip.text = $"提示:复活成功后将会有{CFG_XMLConfig.xmlConfig.GVGInfo.protecttime}秒的无敌保护时间";
                timeStr = "死亡复活倒计时:";
                break;
            case ShowType.MyTeamFirstAttack:
                Count = TurnBasedFightCtrl.Instance.battleWaitTime;
                GetView().txt_tip.text = "我方先手";
                timeStr = "开战倒计时:";
                break;
            case ShowType.OtherTeamFirstAttack:
                Count = TurnBasedFightCtrl.Instance.battleWaitTime;
                GetView().txt_tip.text = "对方先手";
                timeStr = "开战复活倒计时:";
                break;
            default:
                break;
        }

        GameHelper.Instance.AddInvoke(new TimeItem(Count + 1 , 1, () => {
            if (isNullGo) return;
            GetView().txt_Time.text = timeStr + Count;
            if (Count == 0)
                UITool.Close<UIPVPBeginTime>();
            Count--;
        }));
    }

    public override void OnExit()
    {
        GameRoomM.Instance._GameState = GameRoomM.GameState.Play;
    }

    public override void OnDestroy()
    {

    }
}